package util

import (
	"time"

	"github.com/spf13/viper"
)

/*
Stores the configuration for the application
*/
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddr           string        `mapstructure:"SERVER_ADDRESS"`
	SymetricKey          string        `mapstructure:"SYMETRIC_KEY"`
	TokenDuration        time.Duration `mapstructure:"TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

/*
Reads configurations from a configuration file or from
*/
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // json | xml | ... etc

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
