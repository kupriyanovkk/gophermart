package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/kupriyanovkk/gophermart/internal/env"
	"github.com/kupriyanovkk/gophermart/internal/models"
	"github.com/kupriyanovkk/gophermart/internal/store"
	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
)

func Register(w http.ResponseWriter, r *http.Request, s store.Store) {
	var req models.Credentials
	dec := json.NewDecoder(r.Body)
	environ := env.Get()

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, store.ErrorInvalidRequests.Error(), http.StatusBadRequest)
		return
	}

	userID, err := s.RegisterUser(r.Context(), strings.TrimSpace(req.Login), strings.TrimSpace(req.Password))
	if errors.Is(err, store.ErrorLoginConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenutil.SetTokenToHeader(w, userID, environ)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
