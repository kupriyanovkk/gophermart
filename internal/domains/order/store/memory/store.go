package store

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/kupriyanovkk/gophermart/internal/accrual"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/failure"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/models"
)

type MemoryStore struct {
	orders map[int]models.Order
	Chan   chan models.Order
}

var orderStoreSync = sync.RWMutex{}

func (s *MemoryStore) AddOrder(ctx context.Context, orderID, userID int) error {
	orderStoreSync.RLock()
	defer orderStoreSync.RUnlock()

	for _, o := range s.orders {
		if o.Number == strconv.Itoa(orderID) {
			if o.UserID == userID {
				return failure.ErrorOrderAlreadyAdded
			} else {
				return failure.ErrorOrderConflict
			}
		}
	}

	date := time.Now().Format(time.RFC3339)
	s.orders[userID] = models.Order{
		UserID:     userID,
		Number:     strconv.Itoa(orderID),
		Status:     accrual.StatusNew,
		Accrual:    0,
		UploadedAt: date,
	}

	return nil
}

func (s *MemoryStore) UpdateOrder(ctx context.Context, operation accrual.Accrual) error {
	orderStoreSync.RLock()
	defer orderStoreSync.RUnlock()

	for _, o := range s.orders {
		if o.Number == operation.Order {
			order, _ := strconv.Atoi(operation.Order)
			s.orders[order] = models.Order{
				Number:  operation.Order,
				Status:  operation.Status,
				Accrual: operation.Accrual,
			}
			break
		}
	}

	return nil
}

func (s *MemoryStore) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	orderStoreSync.RLock()
	defer orderStoreSync.RUnlock()

	var result []models.Order

	for _, o := range s.orders {
		if o.UserID == userID {
			result = append(result, o)
		}
	}

	return result, nil
}

func (s *MemoryStore) ReadChan() models.Order {
	return <-s.Chan
}

func (s *MemoryStore) WriteChan(order models.Order) {
	s.Chan <- order
}

func NewStore() models.OrderStore {
	return &MemoryStore{
		orders: make(map[int]models.Order),
		Chan:   make(chan models.Order, 100),
	}
}
