package models

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

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}
