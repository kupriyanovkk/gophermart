package flush

import (
	"context"
	"strconv"
	"time"

	"github.com/kupriyanovkk/gophermart/internal/accrual"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/models"
	"go.uber.org/zap"
)

func Run(store models.OrderStore, accrualChan chan accrual.Accrual, client accrual.Client) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	for {
		order := store.ReadChan()
		orderID, _ := strconv.Atoi(order.Number)
		status, err := client.CheckStatus(orderID)

		sugar.Infoln("status", status)

		if err != nil {
			sugar.Errorln(
				"Order Flush",
				"err", err.Error(),
				"orderID", orderID,
			)
			return
		}

		if status.Status != accrual.StatusNotRegister {
			err = store.UpdateOrder(context.TODO(), status)
			if err != nil {
				sugar.Errorln(
					"store.UpdateOrder",
					"err", err.Error(),
					"status", status,
				)
				return
			}

			if status.Status != accrual.StatusProcessed {
				time.AfterFunc(5*time.Second, func() {
					store.WriteChan(models.Order{
						UserID: order.UserID,
						Number: order.Number,
					})
				})
				return
			}

			if status.Accrual > 0 {
				accrualChan <- status
			}
		}
	}
}
