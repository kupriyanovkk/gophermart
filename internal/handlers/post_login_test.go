package handlers

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kupriyanovkk/gophermart/internal/cryptoutil"
	"github.com/kupriyanovkk/gophermart/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	encrypt, _ := cryptoutil.Get()
	store := &store.Store{DB: db, Encrypt: encrypt}
	encryptedPass := store.Encrypt.AEAD.Seal(nil, store.Encrypt.Nonce, []byte("testpassword"), nil)
	requestBody := `{"login": "testuser", "password": "testpassword"}`

	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery("SELECT password, id FROM users WHERE login = ?").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"password", "id"}).AddRow(hex.EncodeToString(encryptedPass), 1))

	rr := httptest.NewRecorder()

	Login(rr, req, *store)

	assert.Equal(t, http.StatusOK, rr.Code)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
