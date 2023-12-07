package models

import (
	"context"

	"github.com/kupriyanovkk/gophermart/internal/accrual"
)

type Good struct {
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type OrderInfo struct {
	ID    int `json:"order"`
	Goods []Good
}

type Order struct {
	UserID     int
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type OrderStore interface {
	ReadChan() Order
	WriteChan(order Order)
	AddOrder(ctx context.Context, orderID int, userID int) error
	UpdateOrder(ctx context.Context, orderStatus accrual.Accrual) error
	GetOrders(ctx context.Context, userID int) ([]Order, error)
}
