package balance

import (
	"context"

	"github.com/kupriyanovkk/gophermart/internal/accrual"
	"github.com/kupriyanovkk/gophermart/internal/domains/balance/models"
	memStore "github.com/kupriyanovkk/gophermart/internal/domains/balance/store/memory"
	pgStore "github.com/kupriyanovkk/gophermart/internal/domains/balance/store/pg"
	"github.com/kupriyanovkk/gophermart/internal/shared"
	"go.uber.org/zap"
)

type Balance struct {
	Store models.BalanceStore
}

func (b *Balance) Flush() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	for {
		operation := b.Store.ReadChan()

		err := b.Store.AddPoints(context.TODO(), operation)

		if err != nil {
			sugar.Errorln(
				"Balance Flush",
				"err", err.Error(),
				"operation", operation,
			)
			return
		}
	}
}

func NewBalance(db shared.DatabaseConnection, accrualChan chan accrual.Accrual) Balance {
	var store models.BalanceStore

	if db == nil {
		store = memStore.NewStore(accrualChan)
	} else {
		store = pgStore.NewStore(db, accrualChan)
	}

	balance := Balance{
		Store: store,
	}

	go balance.Flush()

	return balance
}
