package middlewares

import (
	"errors"
	"net/http"

	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

var ErrorUnauthorized = errors.New("Unauthorized")

func JwtAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/api/user/register" && r.RequestURI != "/api/user/login" {
			userID := tokenutil.GetUserIDFromHeader(r)

			if userID == -1 {
				http.Error(w, ErrorUnauthorized.Error(), http.StatusUnauthorized)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
