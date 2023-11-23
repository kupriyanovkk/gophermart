package shared

import (
	"context"
	"database/sql"

	"github.com/kupriyanovkk/gophermart/internal/cryptoutil"
)

type DatabaseConnection interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Store struct {
	DB      DatabaseConnection
	Encrypt cryptoutil.Encrypt
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
			status VARCHAR(128),
			accrual NUMERIC(10,2),
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
			fk_order_id BIGINT REFERENCES orders(id) NOT NULL,
			date DATE DEFAULT CURRENT_DATE,
			fk_balance_id INTEGER REFERENCES balance(id) NOT NULL
		)
	`)

	return tx.Commit()
}
