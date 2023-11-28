package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kupriyanovkk/gophermart/internal/domains/balance/models"
	"github.com/kupriyanovkk/gophermart/internal/shared"
)

type Store struct {
	db          shared.DatabaseConnection
	LoyaltyChan chan shared.LoyaltyOperation
}

var (
	ErrorNoMoney      = errors.New("there are not enough funds in the account")
	ErrorInvalidOrder = errors.New("invalid order number")
)

func (s *Store) GetUserBalance(ctx context.Context, userID int) (models.UserBalance, error) {
	var (
		current   sql.NullFloat64
		withdrawn sql.NullFloat64
	)
	row := s.db.QueryRowContext(ctx, `SELECT current, withdrawn FROM balance WHERE fk_user_id = $1`, userID)
	err := row.Scan(&current, &withdrawn)

	if err != nil {
		return models.UserBalance{
			Current:   0,
			Withdrawn: 0,
		}, err
	}

	return models.UserBalance{
		Current:   float32(current.Float64),
		Withdrawn: float32(withdrawn.Float64),
	}, nil
}

func (s *Store) AddPoints(ctx context.Context, orderID string, accrual, withdrawn float32) error {
	var (
		userID int
	)
	row := s.db.QueryRowContext(ctx, `SELECT fk_user_id FROM orders WHERE id = $1`, orderID)
	err := row.Scan(&userID)

	if err != nil {
		return err
	}

	userBalance, err := s.GetUserBalance(ctx, userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if userBalance.Current > 0 {
		userBalance.Current += accrual
		return s.UpdateUserBalance(ctx, userID, orderID, userBalance)
	} else {
		userBalance.Current = accrual
		return s.InsertUserBalance(ctx, userID, userBalance.Current)
	}
}

func (s *Store) InsertUserBalance(ctx context.Context, userID int, current float32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO balance (current, withdrawn, fk_user_id)
		VALUES($1, $2, $3);
	`, current, nil, userID)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (s *Store) UpdateUserBalance(ctx context.Context, userID int, orderID string, userBalance models.UserBalance) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE balance SET current = $1, withdrawn = $2
		WHERE fk_user_id = $3
	`, userBalance.Current, userBalance.Withdrawn, userID)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (s *Store) AddWithdraw(ctx context.Context, userID int, orderID string, sum float32) error {
	userBalance, err := s.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}

	userBalance.Current -= sum
	userBalance.Withdrawn += sum

	if userBalance.Current < 0 {
		return ErrorNoMoney
	}

	err = s.UpdateUserBalance(ctx, userID, orderID, userBalance)

	if err == sql.ErrNoRows {
		err = ErrorInvalidOrder
	}

	fmt.Println(err)

	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	date := time.Now().Format(time.RFC3339)
	_, err = tx.ExecContext(ctx, `
		INSERT INTO withdrawals (order_id, date, sum, fk_user_id)
		VALUES($1, $2, $3, $4);
	`, orderID, date, sum, userID)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (s *Store) SelectWithdraws(ctx context.Context, userID int) ([]models.Withdraws, error) {
	limit := 100
	result := make([]models.Withdraws, 0, limit)

	rows, err := s.db.QueryContext(ctx, `SELECT order_id, sum, date FROM withdrawals WHERE fk_user_id = $1 LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var w models.Withdraws
		err = rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, models.Withdraws{
			Order:       w.Order,
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewStore(db shared.DatabaseConnection, loyaltyChan chan shared.LoyaltyOperation) *Store {
	return &Store{db: db, LoyaltyChan: loyaltyChan}
}
