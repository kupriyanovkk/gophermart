package store

import (
	"context"
	"errors"
	"time"

	"github.com/kupriyanovkk/gophermart/internal/domains/order/models"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/status"
	"github.com/kupriyanovkk/gophermart/internal/shared"
)

var (
	ErrorOrderConflict     = errors.New("order has already been uploaded by another user")
	ErrorOrderAlreadyAdded = errors.New("order has already been uploaded by this user")
)

type Store struct {
	db         shared.DatabaseConnection
	OrdersChan chan int
}

func (s *Store) AddOrder(ctx context.Context, orderID, userID int) error {
	var (
		user int
	)
	row := s.db.QueryRowContext(ctx, `SELECT fk_user_id FROM orders WHERE id = $1`, orderID)
	err := row.Scan(&user)

	if err == nil && user != 0 {
		// order already added
		if user == userID {
			err = ErrorOrderAlreadyAdded
		} else {
			err = ErrorOrderConflict
		}

		return err
	}

	date := time.Now().Format(time.RFC3339)
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO orders
		(id, status, accrual, date, fk_user_id)
		VALUES
		($1, $2, $3, $4, $5);
	`, orderID, status.OrderStatusNew, nil, date, userID)

	return err
}

func (s *Store) UpdateOrder(ctx context.Context, orderStatus shared.LoyaltyOperation) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE orders SET status = $1, accrual = $2
			WHERE id = $3
	`, orderStatus.Status, orderStatus.Accrual, orderStatus.ID)

	return err
}

func (s *Store) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	limit := 100
	result := make([]models.Order, 0, limit)

	rows, err := s.db.QueryContext(ctx, `SELECT id, status, accrual, date FROM orders WHERE fk_user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o models.Order
		err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, models.Order{
			Number:     o.Number,
			Status:     o.Status,
			Accrual:    o.Accrual,
			UploadedAt: o.UploadedAt,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewStore(db shared.DatabaseConnection) *Store {
	return &Store{db: db, OrdersChan: make(chan int)}
}
