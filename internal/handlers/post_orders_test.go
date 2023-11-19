package handlers

// import (
// 	"bytes"
// 	"database/sql"
// 	"net/http"
// 	"net/http/httptest"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/kupriyanovkk/gophermart/internal/cryptoutil"
// 	"github.com/kupriyanovkk/gophermart/internal/env"
// 	"github.com/kupriyanovkk/gophermart/internal/order"
// 	"github.com/kupriyanovkk/gophermart/internal/store"
// 	"github.com/kupriyanovkk/gophermart/internal/tokenutil"
// 	"github.com/stretchr/testify/assert"
// )

// func TestPostOrdersHandler(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer db.Close()

// 	encrypt, _ := cryptoutil.Get()
// 	store := &store.Store{DB: db, Encrypt: encrypt}
// 	environ := env.Get()
// 	orderID := 12345678903
// 	requestBody := strconv.Itoa(orderID)
// 	req, err := http.NewRequest("POST", "/orders", bytes.NewBufferString(requestBody))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	token, _ := tokenutil.BuildJWTString(1, environ)
// 	bearer := "Bearer " + token
// 	req.Header.Set("Authorization", bearer)

// 	mock.ExpectQuery("SELECT fk_user_id FROM orders WHERE id = ?").
// 		WithArgs(orderID).
// 		WillReturnError(sql.ErrNoRows)

// 	mock.ExpectExec("INSERT INTO orders").
// 		WithArgs(orderID, order.OrderStatusNew, time.Now().Format(time.RFC3339), 1).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	rr := httptest.NewRecorder()

// 	PostOrders(rr, req, *store, make(chan int, 10))

// 	assert.Equal(t, http.StatusAccepted, rr.Code)

// 	err = mock.ExpectationsWereMet()
// 	if err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }
