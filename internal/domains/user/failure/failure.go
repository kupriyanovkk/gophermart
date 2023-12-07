package failure

import "errors"

var (
	ErrorInvalidCredentials = errors.New("invalid login/password pair")
	ErrorLoginConflict      = errors.New("login is already occupied")
	ErrorInvalidRequests    = errors.New("invalid request format")
)
