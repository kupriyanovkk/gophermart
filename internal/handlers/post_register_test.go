package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kupriyanovkk/gophermart/internal/cryptoutil"
	"github.com/kupriyanovkk/gophermart/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	encrypt, _ := cryptoutil.Get()
	store := &store.Store{DB: db, Encrypt: encrypt}

	requestBody := `{"login": "testuser", "password": "testpassword"}`
	req, err := http.NewRequest("POST", "/register", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	rr := httptest.NewRecorder()

	Register(rr, req, *store)

	assert.Equal(t, http.StatusOK, rr.Code)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
