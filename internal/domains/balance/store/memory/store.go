package store

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/kupriyanovkk/gophermart/internal/accrual"
	"github.com/kupriyanovkk/gophermart/internal/domains/balance/failure"
	"github.com/kupriyanovkk/gophermart/internal/domains/balance/models"
)

type MemoryStore struct {
	balance     map[int]models.UserBalance
	withdrawals []models.Withdraw
	AccrualChan chan accrual.Accrual
}

var balanceStoreSync = sync.RWMutex{}

func (s *MemoryStore) GetUserBalance(ctx context.Context, userID int) (models.UserBalance, error) {
	balanceStoreSync.RLock()
	defer balanceStoreSync.RUnlock()

	if value, ok := s.balance[userID]; ok {
		return models.UserBalance{
			Current:   float32(value.Current),
			Withdrawn: float32(value.Withdrawn),
		}, nil
	}

	return models.UserBalance{
		Current:   0,
		Withdrawn: 0,
	}, sql.ErrNoRows
}

func (s *MemoryStore) AddPoints(ctx context.Context, accrual accrual.Accrual) error {
	userBalance, err := s.GetUserBalance(ctx, accrual.UserID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if userBalance.Current > 0 {
		userBalance.Current += accrual.Accrual
		return s.UpdateUserBalance(ctx, accrual.UserID, accrual.Order, userBalance)
	} else {
		userBalance.Current = accrual.Accrual
		return s.InsertUserBalance(ctx, accrual.UserID, userBalance.Current)
	}
}

func (s *MemoryStore) InsertUserBalance(ctx context.Context, userID int, current float32) error {
	balanceStoreSync.Lock()

	s.balance[userID] = models.UserBalance{
		Current:   float32(current),
		Withdrawn: 0,
	}

	balanceStoreSync.Unlock()

	return nil
}

func (s *MemoryStore) UpdateUserBalance(ctx context.Context, userID int, orderID string, userBalance models.UserBalance) error {
	balanceStoreSync.Lock()

	s.balance[userID] = models.UserBalance{
		Current:   userBalance.Current,
		Withdrawn: userBalance.Withdrawn,
	}

	balanceStoreSync.Unlock()

	return nil
}

func (s *MemoryStore) AddWithdraw(ctx context.Context, userID int, orderID string, sum float32) error {
	userBalance, err := s.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}

	userBalance.Current -= sum
	userBalance.Withdrawn += sum

	if userBalance.Current < 0 {
		return failure.ErrorNoMoney
	}

	// balanceStoreSync.Lock()
	// defer balanceStoreSync.Unlock()
	err = s.UpdateUserBalance(ctx, userID, orderID, userBalance)

	if err == sql.ErrNoRows {
		err = failure.ErrorInvalidOrder
	}

	if err != nil {
		return err
	}

	date := time.Now().Format(time.RFC3339)
	s.withdrawals = append(s.withdrawals, models.Withdraw{
		UserID:      userID,
		Order:       orderID,
		Sum:         sum,
		ProcessedAt: date,
	})

	return nil
}

func (s *MemoryStore) SelectWithdraws(ctx context.Context, userID int) ([]models.Withdraw, error) {
	balanceStoreSync.RLock()
	defer balanceStoreSync.RUnlock()

	result := make([]models.Withdraw, 0, 100)

	for _, w := range s.withdrawals {
		if w.UserID == userID {
			result = append(result, w)
		}
	}

	return result, nil
}

func (s *MemoryStore) ReadChan() accrual.Accrual {
	return <-s.AccrualChan
}

func (s *MemoryStore) WriteChan(accrual accrual.Accrual) {
	s.AccrualChan <- accrual
}

func NewStore(accrualChan chan accrual.Accrual) models.BalanceStore {
	return &MemoryStore{
		AccrualChan: accrualChan,
		balance:     make(map[int]models.UserBalance),
		withdrawals: make([]models.Withdraw, 100),
	}
}
