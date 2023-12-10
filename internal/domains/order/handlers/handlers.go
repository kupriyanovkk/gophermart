package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/kupriyanovkk/gophermart/internal/domains/order/failure"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/models"
	"github.com/kupriyanovkk/gophermart/internal/luhn"
	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

var orderStore models.OrderStore

func Init(store models.OrderStore) {
	orderStore = store
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := tokenutil.GetUserIDFromHeader(r)
	orders, err := orderStore.GetOrders(r.Context(), userID)

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

func PostOrders(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	orderID, err := strconv.Atoi(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !luhn.Valid(orderID) {
		http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	userID := tokenutil.GetUserIDFromHeader(r)
	err = orderStore.AddOrder(r.Context(), orderID, userID)
	if errors.Is(err, failure.ErrorOrderConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err == nil || errors.Is(err, failure.ErrorOrderAlreadyAdded) {
		orderStore.WriteChan(models.Order{
			UserID: userID,
			Number: strconv.Itoa(orderID),
		})
	}

	if err != nil {
		if errors.Is(err, failure.ErrorOrderAlreadyAdded) {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, err.Error(), http.StatusConflict)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
