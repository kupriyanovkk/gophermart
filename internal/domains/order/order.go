package order

import (
	"github.com/kupriyanovkk/gophermart/internal/accrual"
	accrualFlush "github.com/kupriyanovkk/gophermart/internal/accrual/flush"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/models"
	memStore "github.com/kupriyanovkk/gophermart/internal/domains/order/store/memory"
	pgStore "github.com/kupriyanovkk/gophermart/internal/domains/order/store/pg"
	"github.com/kupriyanovkk/gophermart/internal/shared"
)

type Order struct {
	Store models.OrderStore
}

func NewOrder(db shared.DatabaseConnection, accrualChan chan accrual.Accrual, accrualClient accrual.Client) Order {
	var store models.OrderStore

	if db == nil {
		store = memStore.NewStore()
	} else {
		store = pgStore.NewStore(db)
	}

	order := Order{
		Store: store,
	}

	go accrualFlush.Run(order.Store, accrualChan, accrualClient)

	return order
}
