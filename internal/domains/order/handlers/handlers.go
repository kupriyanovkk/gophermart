package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/kupriyanovkk/gophermart/internal/config"
	"github.com/kupriyanovkk/gophermart/internal/domains/order"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/store"
	"github.com/kupriyanovkk/gophermart/internal/luhn"
	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

var orderStore *store.Store

func init() {
	fmt.Println("order init")
	flags := config.Get()
	db, err := sql.Open("postgres", flags.DatabaseURI)

	if err != nil {
		panic(err)
	}

	orderStore = store.NewStore(db)

	go order.Flush(orderStore)
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
	if err != nil || !luhn.Valid(orderID) {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	userID := tokenutil.GetUserIDFromHeader(r)
	err = orderStore.AddOrder(r.Context(), orderID, userID)
	if errors.Is(err, store.ErrorOrderConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err == nil || errors.Is(err, store.ErrorOrderAlreadyAdded) {
		orderStore.OrdersChan <- orderID
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
