package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/kupriyanovkk/gophermart/internal/domains/user/failure"
	"github.com/kupriyanovkk/gophermart/internal/domains/user/models"
	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

var userStore models.UserStore

func Init(store models.UserStore) {
	userStore = store
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.Credentials
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, failure.ErrorInvalidRequests.Error(), http.StatusBadRequest)
		return
	}

	userID, err := userStore.LoginUser(r.Context(), strings.TrimSpace(req.Login), strings.TrimSpace(req.Password))
	if errors.Is(err, failure.ErrorInvalidCredentials) {
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
		http.Error(w, failure.ErrorInvalidRequests.Error(), http.StatusBadRequest)
		return
	}

	userID, err := userStore.RegisterUser(r.Context(), strings.TrimSpace(req.Login), strings.TrimSpace(req.Password))
	if errors.Is(err, failure.ErrorLoginConflict) {
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
