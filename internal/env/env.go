package env

import "os"

func GetToken() string {
	accessTokenSecret, tokenExists := os.LookupEnv("ACCESS_TOKEN_SECRET")

	if tokenExists {
		return accessTokenSecret
	}

	return "access_token_secret"
}

func GetTokenExpiryHour() string {
	accessTokenExpiryHour, tokenExpExists := os.LookupEnv("ACCESS_TOKEN_EXPIRY_HOUR")

	if tokenExpExists {
		return accessTokenExpiryHour
	}

	return "3"
}
