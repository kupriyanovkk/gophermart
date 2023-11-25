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
	orderStatus "github.com/kupriyanovkk/gophermart/internal/domains/order/status"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/store"
	"github.com/kupriyanovkk/gophermart/internal/shared"
	"go.uber.org/zap"
)

func Flush(store *store.Store, loyaltyChan chan shared.LoyaltyOperation) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()
	flags := config.Get()

	for orderID := range store.OrdersChan {
		status, err := CheckStatus(orderID, flags.AccrualSystemAddress)

		sugar.Infoln("status", status)

		if err != nil {
			sugar.Errorln(
				"ORDER Flush",
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
				loyaltyChan <- status
			}
		}
	}
}

func CheckStatus(orderID int, accrualAddr string) (shared.LoyaltyOperation, error) {
	endpoint := fmt.Sprintf("%s/api/orders/%d", accrualAddr, orderID)
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewBuffer(nil))
	if err != nil {
		return shared.LoyaltyOperation{}, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if err != nil {
		return shared.LoyaltyOperation{}, err
	}

	defer response.Body.Close()

	var resp shared.LoyaltyOperation
	if response.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		err := json.Unmarshal(body, &resp)
		if err != nil {
			return shared.LoyaltyOperation{}, err
		}
	}
	if response.StatusCode == http.StatusNoContent {
		return shared.LoyaltyOperation{
			ID:     "-1",
			Status: orderStatus.OrderStatusNotRegister,
		}, nil
	}

	fmt.Println("CheckStatus resp", resp)

	return resp, nil
}
