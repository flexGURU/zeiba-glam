package pkg

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	SERVER_ADDRESS         string        `mapstructure:"SERVER_ADDRESS"`
	DATABASE_URL           string        `mapstructure:"DATABASE_URL"`
	ENVIRONMENT            string        `mapstructure:"ENVIRONMENT"`
	MIGRATION_PATH         string        `mapstructure:"MIGRATION_PATH"`
	PASSWORD_COST          int           `mapstructure:"PASSWORD_COST"`
	REFRESH_TOKEN_DURATION time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TOKEN_DURATION         time.Duration `mapstructure:"TOKEN_DURATION"`
	REDIS_ADDRESS          string        `mapstructure:"REDIS_ADDRESS"`
	REDIS_PASSWORD         string        `mapstructure:"REDIS_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	if err = viper.Unmarshal(&config); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found")
		} else {
			return
		}
	}

	return
}
