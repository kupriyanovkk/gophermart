package failure

import "errors"

var (
	ErrorOrderConflict     = errors.New("order has already been uploaded by another user")
	ErrorOrderAlreadyAdded = errors.New("order has already been uploaded by this user")
)
