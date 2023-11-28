package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kupriyanovkk/gophermart/internal/domains/balance/models"
	"github.com/kupriyanovkk/gophermart/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestGetUserBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	loyaltyChan := make(chan shared.LoyaltyOperation)
	testStore := NewStore(db, loyaltyChan)

	userID := 1
	expectedBalance := models.UserBalance{
		Current:   100.0,
		Withdrawn: 50.0,
	}

	rows := sqlmock.NewRows([]string{"current", "withdrawn"}).
		AddRow(100.0, 50.0)
	mock.ExpectQuery(`^SELECT current, withdrawn FROM balance WHERE fk_user_id`).WithArgs(userID).WillReturnRows(rows)

	balance, err := testStore.GetUserBalance(context.Background(), userID)

	assert.NoError(t, err, "Unexpected error fetching user balance")
	assert.Equal(t, expectedBalance, balance, "Unexpected user balance")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetUserBalance_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	loyaltyChan := make(chan shared.LoyaltyOperation)
	testStore := NewStore(db, loyaltyChan)

	userID := 1

	mock.ExpectQuery(`^SELECT current, withdrawn FROM balance WHERE fk_user_id`).WithArgs(userID).WillReturnError(sql.ErrNoRows)

	balance, err := testStore.GetUserBalance(context.Background(), userID)

	assert.Error(t, err, "Expected error fetching user balance")
	assert.Equal(t, models.UserBalance{Current: 0, Withdrawn: 0}, balance, "Expected default balance on error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestUpdateUserBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	loyaltyChan := make(chan shared.LoyaltyOperation)
	testStore := NewStore(db, loyaltyChan)

	userID := 1
	orderID := "4571398389"
	userBalance := models.UserBalance{
		Current:   100.0,
		Withdrawn: 50.0,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE balance SET current =").WithArgs(userBalance.Current, userBalance.Withdrawn, userID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = testStore.UpdateUserBalance(context.Background(), userID, orderID, userBalance)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("UpdateUserBalance returned an unexpected error: %v", err)
	}
}

func TestInsertUserBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	loyaltyChan := make(chan shared.LoyaltyOperation)
	testStore := NewStore(db, loyaltyChan)

	userID := 1
	var currentBalance float32 = 100.0

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO balance").WithArgs(currentBalance, nil, userID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = testStore.InsertUserBalance(context.Background(), userID, currentBalance)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("InsertUserBalance returned an unexpected error: %v", err)
	}
}

func TestSelectWithdraws(t *testing.T) {
	t.Run("SuccessfulSelectWithdraws", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock: %v", err)
		}
		defer db.Close()

		loyaltyChan := make(chan shared.LoyaltyOperation)
		testStore := NewStore(db, loyaltyChan)

		date := time.Now().Format(time.RFC3339)
		sum := float32(50.0)
		order := "order123"

		rows := sqlmock.NewRows([]string{"order_id", "sum", "date"}).
			AddRow(order, sum, date)

		mock.ExpectQuery("SELECT (.+) FROM withdrawals").
			WithArgs(1, 100).
			WillReturnRows(rows)

		result, err := testStore.SelectWithdraws(context.Background(), 1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedResult := []models.Withdraws{
			{
				Order:       order,
				Sum:         sum,
				ProcessedAt: date,
			},
		}

		if len(result) != len(expectedResult) {
			t.Errorf("Unexpected result length. Expected: %v, Got: %v", len(expectedResult), len(result))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})
}
