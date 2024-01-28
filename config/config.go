package config

import (
	"gotemplate/logger"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AppName            string
	AppEnv             string
	DBConnection       string
	TokenSymmetricKey  string
	HttpUrl            string
	HttpPort           string
	DBHost             string
	DBPort             string
	DBdatabase         string
	DBUsername         string
	DBPassword         string
	TokenDuration      string
	RedisServer        string
	RedisPassword      string
	HttpAllowedOrigins string
	Loglevel           string
	ShutDownTime       string
	ShutDowntype       string
}

func Load() Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	c := Config{}
	log := logger.New("DEBUG")

	// Attempt to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Failed to read configuration", err)
		os.Exit(1)
	}

	// Unmarshal the configuration into the struct
	if err := viper.Unmarshal(&c); err != nil {
		log.Error("Failed to unmarshal configuration", err)
		os.Exit(1)
	}

	// Assign configuration values to environment variables
	os.Setenv("APP_NAME", c.AppName)
	os.Setenv("APP_ENV", c.AppEnv)
	os.Setenv("DB_CONNECTION", c.DBConnection)
	os.Setenv("TOKEN_SYMMETRIC_KEY", c.TokenSymmetricKey)
	os.Setenv("HTTP_URL", c.HttpUrl)
	os.Setenv("HTTP_PORT", c.HttpPort)
	os.Setenv("DB_HOST", c.DBHost)
	os.Setenv("DB_PORT", c.DBPort)
	os.Setenv("DB_DATABASE", c.DBdatabase)
	os.Setenv("DB_USERNAME", c.DBUsername)
	os.Setenv("DB_PASSWORD", c.DBPassword)
	os.Setenv("HTTP_ALLOWED_ORIGINS", c.HttpAllowedOrigins)
	os.Setenv("REDIS_SERVER", c.RedisServer)
	os.Setenv("REDIS_PASSWORD", c.RedisPassword)
	os.Setenv("TOKEN_DURATION", c.TokenDuration)
	os.Setenv("LOG_LEVEL", c.Loglevel)
	os.Setenv("SHUTDOWN_TIME", c.ShutDownTime)
	os.Setenv("SHUTDOWN_TYPE", c.ShutDowntype)

	return c
}
