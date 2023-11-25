package models

type UserBalance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type Withdraw struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type Withdraws struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
