package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kupriyanovkk/gophermart/internal/domains/balance/failure"
	"github.com/kupriyanovkk/gophermart/internal/domains/balance/models"
	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

var balanceStore models.BalanceStore

func Init(store models.BalanceStore) {
	balanceStore = store
}

func GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID := tokenutil.GetUserIDFromHeader(r)
	userBalance, err := balanceStore.GetUserBalance(r.Context(), userID)

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

func PostWithdraw(w http.ResponseWriter, r *http.Request) {
	var req models.Withdraw
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := tokenutil.GetUserIDFromHeader(r)

	err := balanceStore.AddWithdraw(r.Context(), userID, req.Order, req.Sum)
	if errors.Is(err, failure.ErrorNoMoney) {
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func GetWithdraw(w http.ResponseWriter, r *http.Request) {
	userID := tokenutil.GetUserIDFromHeader(r)
	withdraws, err := balanceStore.SelectWithdraws(r.Context(), userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(withdraws) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(withdraws); err != nil {
		return
	}
}
