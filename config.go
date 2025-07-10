package main

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

type AppConfig struct {
	AppID          int64
	PrivateKeyPath string
	Environment    string
	RedisHost      string
	RedisPort      int
	RedisUsername  string
	RedisPassword  string
}

func NewAppConfig() (*AppConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &AppConfig{
		AppID:          viper.GetInt64("APP_ID"),
		PrivateKeyPath: viper.GetString("PRIVATE_KEY_PATH"),
		Environment:    viper.GetString("ENVIRONMENT"),
		RedisHost:      viper.GetString("REDIS_HOST"),
		RedisPort:      viper.GetInt("REDIS_PORT"),
		RedisUsername:  viper.GetString("REDIS_USERNAME"),
		RedisPassword:  viper.GetString("REDIS_PASSWORD"),
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *AppConfig) Validate() error {
	return v.ValidateStruct(c,
		v.Field(&c.AppID, v.Required),
		v.Field(&c.PrivateKeyPath, v.Required),
		v.Field(&c.Environment, v.Required),
		v.Field(&c.RedisHost, v.Required),
		v.Field(&c.RedisPort, v.Required, v.Min(1), v.Max(65535)),
		v.Field(&c.RedisUsername, v.Required),
		v.Field(&c.RedisPassword, v.Required),
	)
}
