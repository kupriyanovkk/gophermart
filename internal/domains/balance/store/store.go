package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/kupriyanovkk/gophermart/internal/domains/balance/models"
	"github.com/kupriyanovkk/gophermart/internal/shared"
	"github.com/lib/pq"
)

type Store struct {
	db shared.DatabaseConnection
}

func (s *Store) GetUserBalance(ctx context.Context, userID int) (models.UserBalance, error) {
	var (
		current   sql.NullFloat64
		withdrawn sql.NullFloat64
	)
	row := s.db.QueryRowContext(ctx, `SELECT current, withdrawn FROM balance WHERE fk_user_id = $1`, userID)
	err := row.Scan(&current, &withdrawn)

	if err != nil {
		return models.UserBalance{}, err
	}

	return models.UserBalance{
		Current:   float32(current.Float64),
		Withdrawn: float32(withdrawn.Float64),
	}, nil
}

func (s *Store) UpdateUserBalance(ctx context.Context, orderID string, accrual float32) error {
	var (
		userID  int
		current float32
	)
	row := s.db.QueryRowContext(ctx, `SELECT fk_user_id FROM orders WHERE id = $1`, orderID)
	err := row.Scan(&userID)

	if err != nil {
		return err
	}

	row = s.db.QueryRowContext(ctx, `SELECT current FROM balance WHERE fk_user_id = $1`, userID)
	err = row.Scan(&current)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code != pgerrcode.NoData {
			return err
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if current > 0 {
		_, err = tx.ExecContext(ctx, `
				UPDATE balance SET current = $1
					WHERE fk_user_id = $2
			`, current+accrual, userID)
	} else {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO balance (current, withdrawn, fk_user_id)
			VALUES($1, $2, $3);
		`, accrual, nil, userID)
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func NewStore(db shared.DatabaseConnection) *Store {
	return &Store{db: db}
}
