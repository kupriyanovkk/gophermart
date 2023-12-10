package tokenutil

import (
	"testing"

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
			tokenString, err := BuildJWTString(tt.userID)

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
			token, _ := BuildJWTString(tt.userID)
			userID := GetUserID(token)

			assert.Equal(t, tt.userID, userID)
		})
	}
}
