package env

import "os"

type Env struct {
	AccessTokenExpiryHour string
	AccessTokenSecret     string
}

func Get() Env {
	accessTokenSecret, tokenExists := os.LookupEnv("ACCESS_TOKEN_SECRET")
	accessTokenExpiryHour, tokenExpExists := os.LookupEnv("ACCESS_TOKEN_EXPIRY_HOUR")
	env := Env{
		AccessTokenSecret:     "",
		AccessTokenExpiryHour: "",
	}

	if tokenExists {
		env.AccessTokenSecret = accessTokenSecret
	}
	if tokenExpExists {
		env.AccessTokenExpiryHour = accessTokenExpiryHour
	}

	return env
}
