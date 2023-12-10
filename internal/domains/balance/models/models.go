package models

import (
	"context"

	"github.com/kupriyanovkk/gophermart/internal/accrual"
)

type BalanceStore interface {
	ReadChan() accrual.Accrual
	WriteChan(accrual accrual.Accrual)
	GetUserBalance(ctx context.Context, userID int) (UserBalance, error)
	AddPoints(ctx context.Context, accrual accrual.Accrual) error
	InsertUserBalance(ctx context.Context, userID int, current float32) error
	UpdateUserBalance(ctx context.Context, userID int, orderID string, userBalance UserBalance) error
	AddWithdraw(ctx context.Context, userID int, orderID string, sum float32) error
	SelectWithdraws(ctx context.Context, userID int) ([]Withdraw, error)
}

type UserBalance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type Withdraw struct {
	UserID      int
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
