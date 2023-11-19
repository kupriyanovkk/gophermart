package models

import (
	"context"
	"database/sql"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kupriyanovkk/gophermart/internal/store"
)

type ConfigFlags struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

type DatabaseConnection interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

type App struct {
	Flags      ConfigFlags
	Store      store.Store
	OrdersChan chan int
}
