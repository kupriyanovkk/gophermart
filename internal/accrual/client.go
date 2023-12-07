package accrual

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	CheckStatus(orderID int) (Accrual, error)
}

type client struct {
	baseURL string
}

func (c client) CheckStatus(orderID int) (Accrual, error) {
	endpoint := fmt.Sprintf("%s/api/orders/%d", c.baseURL, orderID)
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewBuffer(nil))
	if err != nil {
		return Accrual{}, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if err != nil {
		return Accrual{}, err
	}

	defer response.Body.Close()

	var resp Accrual
	if response.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		err := json.Unmarshal(body, &resp)
		if err != nil {
			return Accrual{}, err
		}
	}
	if response.StatusCode == http.StatusNoContent {
		return Accrual{
			Order:  "-1",
			Status: StatusNotRegister,
		}, nil
	}

	fmt.Println("CheckStatus resp", resp)

	return resp, nil
}

func NewClient(baseURL string) Client {
	return client{baseURL}
}
