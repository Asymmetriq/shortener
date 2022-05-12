package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config interface {
	GetAddress() string
	GetBaseURL() string
}

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	ServerBaseURL string `env:"SERVER_BASE_URL"`
}

func (c *config) GetAddress() string {
	if len(c.ServerAddress) == 0 {
		return ":8080"
	}
	return c.ServerAddress
}

func (c *config) GetBaseURL() string {
	return c.ServerBaseURL
}

func InitConfig() Config {
	var conf config
	err := env.Parse(&conf)
	if err != nil {
		log.Fatalf("missing required env vars: %v", err)
	}
	return &conf
}
