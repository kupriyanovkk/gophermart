package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kupriyanovkk/gophermart/internal/domains/order/status"
	"github.com/stretchr/testify/assert"
)

func TestAddOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	testStore := NewStore(db)

	orderID := 1
	userID := 2

	mock.ExpectQuery(`^SELECT fk_user_id FROM orders WHERE id`).WithArgs(orderID).WillReturnRows(sqlmock.NewRows([]string{"fk_user_id"}).AddRow(userID))

	err = testStore.AddOrder(context.Background(), orderID, userID)

	assert.EqualError(t, err, ErrorOrderAlreadyAdded.Error(), "Expected error for order already added")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAddOrder_OrderConflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	testStore := NewStore(db)

	orderID := 1
	userID := 2

	mock.ExpectQuery(`^SELECT fk_user_id FROM orders WHERE id`).WithArgs(orderID).WillReturnRows(sqlmock.NewRows([]string{"fk_user_id"}).AddRow(3))

	err = testStore.AddOrder(context.Background(), orderID, userID)

	assert.EqualError(t, err, ErrorOrderConflict.Error(), "Expected error for order conflict")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAddOrder_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	testStore := NewStore(db)

	orderID := 1
	userID := 2

	mock.ExpectQuery(`^SELECT fk_user_id FROM orders WHERE id`).WithArgs(orderID).WillReturnRows(sqlmock.NewRows([]string{"fk_user_id"}))

	mock.ExpectExec(`^INSERT INTO orders`).WithArgs(orderID, status.OrderStatusNew, 0, sqlmock.AnyArg(), userID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = testStore.AddOrder(context.Background(), orderID, userID)

	assert.NoError(t, err, "Unexpected error adding order")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	testStore := NewStore(db)

	userID := 1

	rows := sqlmock.NewRows([]string{"id", "status", "accrual", "date"}).
		AddRow(1, "Pending", 10, time.Now()).
		AddRow(2, "Completed", 20, time.Now())
	mock.ExpectQuery(`^SELECT id, status, accrual, date FROM orders WHERE fk_user_id`).WithArgs(userID).WillReturnRows(rows)

	orders, err := testStore.GetOrders(context.Background(), userID)

	assert.NoError(t, err, "Unexpected error fetching orders")
	assert.Len(t, orders, 2, "Unexpected number of orders")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetOrders_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	testStore := NewStore(db)

	userID := 1

	mock.ExpectQuery(`^SELECT id, status, accrual, date FROM orders WHERE fk_user_id`).WithArgs(userID).WillReturnError(sql.ErrNoRows)

	orders, err := testStore.GetOrders(context.Background(), userID)

	assert.Error(t, err, "Expected error fetching orders")
	assert.Nil(t, orders, "Orders should be nil on error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
