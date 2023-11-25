package shared

import (
	"context"
	"database/sql"
)

type DatabaseConnection interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type LoyaltyOperation struct {
	ID      string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}

func BootstrapDB(ctx context.Context, DB DatabaseConnection) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			login VARCHAR(128),
			password VARCHAR(128)
		)
	`)

	tx.ExecContext(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS user_login ON users (login)")

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS orders(
			id BIGINT PRIMARY KEY,
			status VARCHAR(128) NOT NULL,
			accrual NUMERIC(10,2) NOT NULL,
			date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_DATE,
			fk_user_id INTEGER REFERENCES users(id) NOT NULL
		)
	`)

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS balance(
			id SERIAL PRIMARY KEY,
			current NUMERIC(10,2),
			withdrawn NUMERIC(10,2),
			fk_user_id INTEGER REFERENCES users(id) NOT NULL
		)
	`)

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS withdrawals(
			order_id BIGINT NOT NULL,
			date DATE DEFAULT CURRENT_DATE,
			sum NUMERIC(10,2),
			fk_user_id INTEGER REFERENCES users(id) NOT NULL
		)
	`)

	return tx.Commit()
}
