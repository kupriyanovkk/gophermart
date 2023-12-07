package store

import (
	"context"
	"database/sql"
	"encoding/hex"
	"sync"

	"github.com/kupriyanovkk/gophermart/internal/cryptoutil"
	"github.com/kupriyanovkk/gophermart/internal/domains/user/failure"
	"github.com/kupriyanovkk/gophermart/internal/domains/user/models"
)

var uuid = 0
var userStorageSync = sync.RWMutex{}

type MemoryStore struct {
	users []models.Credentials
}

func (s *MemoryStore) RegisterUser(ctx context.Context, login, password string) (int, error) {
	userStorageSync.RLock()
	defer userStorageSync.RUnlock()

	for _, u := range s.users {
		if u.Login == login {
			return -1, failure.ErrorLoginConflict
		}
	}

	uuid += 1
	userID := uuid
	encrypt, _ := cryptoutil.Get()
	encryptedPass := encrypt.AEAD.Seal(nil, encrypt.Nonce, []byte(password), nil)

	s.users = append(s.users, models.Credentials{
		ID:       userID,
		Login:    login,
		Password: hex.EncodeToString(encryptedPass),
	})

	return userID, nil
}

func (s *MemoryStore) LoginUser(ctx context.Context, login, password string) (int, error) {
	encrypt, _ := cryptoutil.Get()
	encryptedPass := encrypt.AEAD.Seal(nil, encrypt.Nonce, []byte(password), nil)

	userStorageSync.RLock()
	defer userStorageSync.RUnlock()

	for _, u := range s.users {
		if u.Login == login && hex.EncodeToString(encryptedPass) == u.Password {
			return u.ID, nil
		}
	}

	return -1, failure.ErrorInvalidCredentials
}

func (s *MemoryStore) GetUser(ctx context.Context, userID int) (models.Credentials, error) {
	userStorageSync.RLock()
	defer userStorageSync.RUnlock()

	for _, u := range s.users {
		if u.ID == userID {
			return u, nil
		}
	}

	return models.Credentials{}, sql.ErrNoRows
}

func NewStore() models.UserStore {
	return &MemoryStore{users: make([]models.Credentials, 100)}
}
