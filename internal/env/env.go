package env

import "os"

func GetToken() string {
	accessTokenSecret, tokenExists := os.LookupEnv("ACCESS_TOKEN_SECRET")

	if tokenExists {
		return accessTokenSecret
	}

	return ""
}

func GetTokenExpiryHour() string {
	accessTokenExpiryHour, tokenExpExists := os.LookupEnv("ACCESS_TOKEN_EXPIRY_HOUR")

	if tokenExpExists {
		return accessTokenExpiryHour
	}

	return ""
}
