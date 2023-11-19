package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kupriyanovkk/gophermart/internal/store"
	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

func GetUserBalance(w http.ResponseWriter, r *http.Request, s store.Store) {
	userID := tokenutil.GetUserIDFromHeader(r)
	userBalance, err := s.GetUserBalance(r.Context(), userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(userBalance); err != nil {
		return
	}
}
