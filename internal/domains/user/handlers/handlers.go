package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/kupriyanovkk/gophermart/internal/domains/user/models"
	"github.com/kupriyanovkk/gophermart/internal/domains/user/store"
	"github.com/kupriyanovkk/gophermart/internal/shared"
	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

var userStore *store.Store

func Init(db shared.DatabaseConnection) {
	fmt.Println("user init")

	userStore = store.NewStore(db)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.Credentials
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, store.ErrorInvalidRequests.Error(), http.StatusBadRequest)
		return
	}

	userID, err := userStore.LoginUser(r.Context(), strings.TrimSpace(req.Login), strings.TrimSpace(req.Password))
	if errors.Is(err, store.ErrorInvalidCredentials) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenutil.SetTokenToHeader(w, userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.Credentials
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, store.ErrorInvalidRequests.Error(), http.StatusBadRequest)
		return
	}

	userID, err := userStore.RegisterUser(r.Context(), strings.TrimSpace(req.Login), strings.TrimSpace(req.Password))
	if errors.Is(err, store.ErrorLoginConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenutil.SetTokenToHeader(w, userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
