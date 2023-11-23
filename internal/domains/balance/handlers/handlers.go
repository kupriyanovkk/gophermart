package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kupriyanovkk/gophermart/internal/config"
	"github.com/kupriyanovkk/gophermart/internal/domains/balance/store"
	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

var balanceStore *store.Store

func init() {
	fmt.Println("balance init")
	flags := config.Get()
	db, err := sql.Open("postgres", flags.DatabaseURI)

	if err != nil {
		panic(err)
	}

	balanceStore = store.NewStore(db)
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
