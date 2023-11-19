package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kupriyanovkk/gophermart/internal/store"
	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

func GetOrders(w http.ResponseWriter, r *http.Request, s store.Store) {
	userID := tokenutil.GetUserIDFromHeader(r)
	orders, err := s.GetOrders(r.Context(), userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(orders); err != nil {
		return
	}
}
