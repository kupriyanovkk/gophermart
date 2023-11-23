package store

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/kupriyanovkk/gophermart/internal/cryptoutil"
	"github.com/kupriyanovkk/gophermart/internal/shared"
	"github.com/lib/pq"
)

var (
	ErrorInvalidCredentials = errors.New("invalid login/password pair")
	ErrorLoginConflict      = errors.New("login is already occupied")
	ErrorInvalidRequests    = errors.New("invalid request format")
)

type Store struct {
	db shared.DatabaseConnection
}

func (s *Store) RegisterUser(ctx context.Context, login, password string) (int, error) {
	var userID int
	encrypt, _ := cryptoutil.Get()
	encryptedPass := encrypt.AEAD.Seal(nil, encrypt.Nonce, []byte(password), nil)
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO users
		(login, password)
		VALUES
		($1, $2)
		RETURNING id;
	`, login, hex.EncodeToString(encryptedPass)).Scan(&userID)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = ErrorLoginConflict
		}
	}

	return userID, err
}

func (s *Store) LoginUser(ctx context.Context, login, password string) (int, error) {
	var (
		pass   string
		userID int
	)
	encrypt, _ := cryptoutil.Get()
	encryptedPass := encrypt.AEAD.Seal(nil, encrypt.Nonce, []byte(password), nil)
	row := s.db.QueryRowContext(ctx, `SELECT password, id FROM users WHERE login = $1`, login)
	err := row.Scan(&pass, &userID)

	if err != nil {
		return -1, err
	}

	if hex.EncodeToString(encryptedPass) != pass {
		return -1, ErrorInvalidCredentials
	}

	return userID, nil
}

func NewStore(db shared.DatabaseConnection) *Store {
	return &Store{db: db}
}
