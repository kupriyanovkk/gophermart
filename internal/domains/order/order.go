package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kupriyanovkk/gophermart/internal/config"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/models"
	orderStatus "github.com/kupriyanovkk/gophermart/internal/domains/order/status"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/store"
	"go.uber.org/zap"
)

func Flush(store *store.Store) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()
	flags := config.Get()

	for orderID := range store.OrdersChan {
		status, err := CheckStatus(orderID, flags.AccrualSystemAddress)

		sugar.Infoln("status", status)

		if err != nil {
			sugar.Errorln(
				"err", err.Error(),
				"orderID", orderID,
			)
			return
		}

		if status.Status != orderStatus.OrderStatusNotRegister {
			err = store.UpdateOrder(context.TODO(), status)
			if err != nil {
				sugar.Errorln(
					"err", err.Error(),
					"status", status,
				)
				return
			}

			if status.Status != orderStatus.OrderStatusProcessed {
				time.AfterFunc(5*time.Second, func() {
					store.OrdersChan <- orderID
				})
				return
			}

			if status.Accrual > 0 {
				err := store.UpdateUserBalance(context.TODO(), status.ID, status.Accrual)

				if err != nil {
					sugar.Errorln(
						"err", err.Error(),
						"status", status,
					)
					return
				}
			}
		}
	}
}

func CheckStatus(orderID int, accrualAddr string) (models.OrderAccrual, error) {
	endpoint := fmt.Sprintf("%s/api/orders/%d", accrualAddr, orderID)
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewBuffer(nil))
	if err != nil {
		return models.OrderAccrual{}, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if err != nil {
		return models.OrderAccrual{}, err
	}

	defer response.Body.Close()

	var resp models.OrderAccrual
	if response.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		err := json.Unmarshal(body, &resp)
		if err != nil {
			return models.OrderAccrual{}, err
		}
	}
	if response.StatusCode == http.StatusNoContent {
		return models.OrderAccrual{
			ID:     "-1",
			Status: orderStatus.OrderStatusNotRegister,
		}, nil
	}

	fmt.Println("CheckStatus resp", resp)

	return resp, nil
}
