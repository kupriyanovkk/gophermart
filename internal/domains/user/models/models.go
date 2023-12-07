package models

import "context"

type Credentials struct {
	ID       int
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserStore interface {
	RegisterUser(ctx context.Context, login string, password string) (int, error)
	LoginUser(ctx context.Context, login string, password string) (int, error)
	GetUser(ctx context.Context, userID int) (Credentials, error)
}
