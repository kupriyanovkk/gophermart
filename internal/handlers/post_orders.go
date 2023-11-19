package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/kupriyanovkk/gophermart/internal/luhn"
	"github.com/kupriyanovkk/gophermart/internal/store"
	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

func PostOrders(w http.ResponseWriter, r *http.Request, s store.Store, ordersChan chan int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	orderID, err := strconv.Atoi(string(body))
	if err != nil || !luhn.Valid(orderID) {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	userID := tokenutil.GetUserIDFromHeader(r)
	err = s.AddOrder(r.Context(), orderID, userID)
	if errors.Is(err, store.ErrorOrderConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err == nil || errors.Is(err, store.ErrorOrderAlreadyAdded) {
		ordersChan <- orderID
	}

	if err != nil {
		if errors.Is(err, store.ErrorOrderAlreadyAdded) {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, err.Error(), http.StatusConflict)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
