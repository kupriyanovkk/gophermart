package env

import (
	"github.com/spf13/viper"
)

type Env struct {
	AccessTokenExpiryHour  string
	AccessTokenSecret      string
}

func Get() Env {
	viper.SetConfigFile("../../.env")
	viper.ReadInConfig()

	return Env{
		AccessTokenSecret:     viper.Get("ACCESS_TOKEN_SECRET").(string),
		AccessTokenExpiryHour: viper.Get("ACCESS_TOKEN_EXPIRY_HOUR").(string),
	}
}
