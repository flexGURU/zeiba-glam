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
	FRONTEND_URL           []string      `mapstructure:"FRONTEND_URL"`
	MIGRATION_PATH         string        `mapstructure:"MIGRATION_PATH"`
	PASSWORD_COST          int           `mapstructure:"PASSWORD_COST"`
	TOKEN_SYMMETRIC_KEY    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	REFRESH_TOKEN_DURATION time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	ACCESS_TOKEN_DURATION  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	REDIS_ADDRESS          string        `mapstructure:"REDIS_ADDRESS"`
	REDIS_PASSWORD         string        `mapstructure:"REDIS_PASSWORD"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	setDefaults()

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables")
		} else {
			return Config{}, Errorf(INTERNAL_ERROR, "failed to read config: %s", err.Error())
		}
	}

	var config Config

	return config, viper.Unmarshal(&config)
}

func setDefaults() {
	viper.SetDefault("SERVER_ADDRESS", "")
	viper.SetDefault("DATABASE_URL", "")
	viper.SetDefault("ENVIRONMENT", "")
	viper.SetDefault("FRONTEND_URL", "")
	viper.SetDefault("MIGRATION_PATH", "")
	viper.SetDefault("PASSWORD_COST", 0)
	viper.SetDefault("TOKEN_SYMMETRIC_KEY", "")
	viper.SetDefault("REFRESH_TOKEN_DURATION", 0)
	viper.SetDefault("ACCESS_TOKEN_DURATION", 0)
	viper.SetDefault("REDIS_ADDRESS", "")
	viper.SetDefault("REDIS_PASSWORD", "")
}
