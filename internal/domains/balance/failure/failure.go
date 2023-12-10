package failure

import "errors"

var (
	ErrorNoMoney      = errors.New("there are not enough funds in the account")
	ErrorInvalidOrder = errors.New("invalid order number")
)
