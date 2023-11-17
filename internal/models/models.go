package models

import "github.com/golang-jwt/jwt/v4"

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}
