package tokenutil

import (
	"testing"

	"github.com/kupriyanovkk/gophermart/internal/env"
	"github.com/stretchr/testify/assert"
)

func TestBuildJWTString(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		expiryHour    string
		secret        string
		expectedError bool
	}{
		{
			name:          "ValidToken",
			userID:        123,
			expiryHour:    "1",
			secret:        "mySecret",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := env.Env{
				AccessTokenExpiryHour: tt.expiryHour,
				AccessTokenSecret:     tt.secret,
			}

			tokenString, err := BuildJWTString(tt.userID, env)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name          string
		hour          string
		secret        string
		userID        int
		expectedError bool
	}{
		{
			name:          "ValidCase",
			secret:        "mySecret",
			hour:          "1",
			userID:        123,
			expectedError: false,
		},
		{
			name:          "InvalidCase",
			hour:          "",
			secret:        "",
			userID:        -1,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := env.Env{
				AccessTokenExpiryHour: tt.hour,
				AccessTokenSecret:     tt.secret,
			}

			token, _ := BuildJWTString(tt.userID, env)
			userID := GetUserID(token, env)

			assert.Equal(t, tt.userID, userID)
		})
	}
}
