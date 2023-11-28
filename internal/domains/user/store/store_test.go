package store

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kupriyanovkk/gophermart/internal/cryptoutil"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	store := NewStore(db)

	login := "testuser"
	password := "testpassword"

	mock.ExpectQuery(`^INSERT INTO users`).WithArgs(login, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	userID, err := store.RegisterUser(context.Background(), login, password)

	assert.NoError(t, err, "Unexpected error during user registration")
	assert.Equal(t, 1, userID, "Unexpected user ID")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRegisterUser_UniqueViolationError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	store := NewStore(db)

	login := "existinguser"
	password := "testpassword"

	mock.ExpectQuery(`^INSERT INTO users`).WithArgs(login, sqlmock.AnyArg()).WillReturnError(ErrorLoginConflict)

	userID, err := store.RegisterUser(context.Background(), login, password)

	assert.EqualError(t, err, ErrorLoginConflict.Error(), "Expected unique violation error")
	assert.Zero(t, userID, "User ID should be zero on error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLoginUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	store := NewStore(db)

	login := "testuser"
	password := "testpassword"

	encrypt, _ := cryptoutil.Get()
	encryptedPass := encrypt.AEAD.Seal(nil, encrypt.Nonce, []byte(password), nil)
	mock.ExpectQuery(`^SELECT password, id FROM users WHERE login`).WithArgs(login).WillReturnRows(sqlmock.NewRows([]string{"password", "id"}).AddRow(hex.EncodeToString(encryptedPass), 1))

	userID, err := store.LoginUser(context.Background(), login, password)

	assert.NoError(t, err, "Unexpected error during user login")
	assert.Equal(t, 1, userID, "Unexpected user ID")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating SQL mock: %v", err)
	}
	defer db.Close()

	store := NewStore(db)

	login := "testuser"
	password := "invalidpassword"

	encrypt, _ := cryptoutil.Get()
	encryptedPass := encrypt.AEAD.Seal(nil, encrypt.Nonce, []byte("correctpassword"), nil)
	mock.ExpectQuery(`^SELECT password, id FROM users WHERE login`).WithArgs(login).WillReturnRows(sqlmock.NewRows([]string{"password", "id"}).AddRow(hex.EncodeToString(encryptedPass), 1))

	userID, err := store.LoginUser(context.Background(), login, password)

	assert.EqualError(t, err, ErrorInvalidCredentials.Error(), "Expected invalid credentials error")
	assert.Equal(t, -1, userID, "User ID should be -1 on error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
