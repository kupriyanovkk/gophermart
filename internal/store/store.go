package store

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/kupriyanovkk/gophermart/internal/cryptoutil"
	"github.com/kupriyanovkk/gophermart/internal/order"
	"github.com/lib/pq"
)

var (
	ErrorInvalidCredentials = errors.New("invalid login/password pair")
	ErrorLoginConflict      = errors.New("login is already occupied")
	ErrorInvalidRequests    = errors.New("invalid request format")
	ErrorOrderConflict      = errors.New("order has already been uploaded by another user")
	ErrorOrderAlreadyAdded  = errors.New("order has already been uploaded by this user")
)

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type DatabaseConnection interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Store struct {
	DB      DatabaseConnection
	Encrypt cryptoutil.Encrypt
}

func (s *Store) bootstrap(ctx context.Context) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			login VARCHAR(128),
			password VARCHAR(128)
		)
	`)

	tx.ExecContext(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS user_login ON users (login)")

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS orders(
			id BIGINT PRIMARY KEY,
			status VARCHAR(128),
			accrual NUMERIC(5,2),
			date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_DATE,
			fk_user_id INTEGER REFERENCES users(id) NOT NULL
		)
	`)

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS balance(
			id SERIAL PRIMARY KEY,
			amount NUMERIC(5,2),
			fk_user_id INTEGER REFERENCES users(id) NOT NULL
		)
	`)

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS withdrawals(
			fk_order_id BIGINT REFERENCES orders(id) NOT NULL,
			date DATE DEFAULT CURRENT_DATE,
			fk_balance_id INTEGER REFERENCES balance(id) NOT NULL
		)
	`)

	return tx.Commit()
}

func (s *Store) RegisterUser(ctx context.Context, login, password string) (int, error) {
	var userID int
	encryptedPass := s.Encrypt.AEAD.Seal(nil, s.Encrypt.Nonce, []byte(password), nil)
	err := s.DB.QueryRowContext(ctx, `
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
	encryptedPass := s.Encrypt.AEAD.Seal(nil, s.Encrypt.Nonce, []byte(password), nil)
	row := s.DB.QueryRowContext(ctx, `SELECT password, id FROM users WHERE login = $1`, login)
	err := row.Scan(&pass, &userID)

	if err != nil {
		return -1, err
	}

	if hex.EncodeToString(encryptedPass) != pass {
		return -1, ErrorInvalidCredentials
	}

	return userID, nil
}

func (s *Store) AddOrder(ctx context.Context, orderID, userID int) error {
	var (
		user int
	)
	row := s.DB.QueryRowContext(ctx, `SELECT fk_user_id FROM orders WHERE id = $1`, orderID)
	err := row.Scan(&user)

	if err == nil && user != 0 {
		// order already added
		if user == userID {
			err = ErrorOrderAlreadyAdded
		} else {
			err = ErrorOrderConflict
		}

		return err
	}

	date := time.Now().Format(time.RFC3339)
	_, err = s.DB.ExecContext(ctx, `
		INSERT INTO orders
		(id, status, accrual, date, fk_user_id)
		VALUES
		($1, $2, $3, $4, $5);
	`, orderID, order.OrderStatusNew, nil, date, userID)

	return err
}

func (s *Store) UpdateOrder(ctx context.Context, orderStatus order.OrderAccrual) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE orders SET status = $1, accrual = $2
			WHERE id = $3
	`, orderStatus.Status, orderStatus.Accrual, orderStatus.ID)

	return err
}

func (s *Store) GetOrders(ctx context.Context, userID int) ([]Order, error) {
	limit := 100
	result := make([]Order, 0, limit)

	rows, err := s.DB.QueryContext(ctx, `SELECT id, status, accrual, date FROM orders WHERE fk_user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, Order{
			Number:     o.Number,
			Status:     o.Status,
			Accrual:    o.Accrual,
			UploadedAt: o.UploadedAt,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewStore(db DatabaseConnection) *Store {
	encrypt, _ := cryptoutil.Get()
	store := &Store{DB: db, Encrypt: encrypt}
	err := store.bootstrap(context.Background())

	if err != nil {
		panic(err)
	}

	return store
}
