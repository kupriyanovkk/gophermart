package balance

import (
	"context"

	"github.com/kupriyanovkk/gophermart/internal/domains/balance/store"
	"go.uber.org/zap"
)

func Flush(store *store.Store) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	for operation := range store.LoyaltyChan {
		err := store.AddPoints(context.TODO(), operation.ID, operation.Accrual, 0)

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
