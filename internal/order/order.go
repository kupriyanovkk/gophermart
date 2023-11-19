package order

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Good struct {
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type OrderInfo struct {
	ID    int `json:"order"`
	Goods []Good
}

type OrderAccrual struct {
	ID      string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}

const (
	OrderStatusNew         string = "NEW"
	OrderStatusProcessing  string = "PROCESSING"
	OrderStatusInvalid     string = "INVALID"
	OrderStatusProcessed   string = "PROCESSED"
	OrderStatusNotRegister string = "NOT_REGISTER"
)

func CheckStatus(orderID int, accrualAddr string) (OrderAccrual, error) {
	endpoint := fmt.Sprintf("%s/api/orders/%d", accrualAddr, orderID)
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewBuffer(nil))
	if err != nil {
		return OrderAccrual{}, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if err != nil {
		return OrderAccrual{}, err
	}

	defer response.Body.Close()

	var resp OrderAccrual
	if response.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		err := json.Unmarshal(body, &resp)
		if err != nil {
			return OrderAccrual{}, err
		}
	}
	if response.StatusCode == http.StatusNoContent {
		return OrderAccrual{
			ID:     "-1",
			Status: OrderStatusNotRegister,
		}, nil
	}

	fmt.Println("CheckStatus resp", resp)

	return resp, nil
}
