package tokenutil

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/kupriyanovkk/gophermart/internal/env"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func BuildJWTString(userID int) (string, error) {
	duration, _ := strconv.ParseInt(env.GetTokenExpiryHour(), 10, 64)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(duration))),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(env.GetToken()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(env.GetToken()), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return -1
	}

	return claims.UserID
}

func SetTokenToHeader(w http.ResponseWriter, userID int) {
	token, _ := BuildJWTString(userID)
	bearer := "Bearer " + token

	w.Header().Set("Authorization", bearer)
}

func GetUserIDFromHeader(r *http.Request) int {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return -1
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return -1
	}

	token := strings.TrimSpace(splitToken[1])

	return GetUserID(token)
}
