package accrual

const (
	StatusNew         string = "NEW"
	StatusProcessing  string = "PROCESSING"
	StatusInvalid     string = "INVALID"
	StatusProcessed   string = "PROCESSED"
	StatusNotRegister string = "NOT_REGISTER"
)

type Accrual struct {
	UserID  int
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
