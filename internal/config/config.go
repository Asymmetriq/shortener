package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	address     *string
	baseURL     *string
	storagepath *string
)

func init() {
	address = flag.String("a", ":8080", "Server's host address")
	baseURL = flag.String("b", "", "Server's base URL")
	storagepath = flag.String("f", "", "Storage path")
}

type Config interface {
	GetAddress() string
	GetBaseURL() string
	GetStoragePath() string
}

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	ServerBaseURL   string `env:"SERVER_BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func (c *config) GetAddress() string {
	return c.ServerAddress
}

func (c *config) GetBaseURL() string {
	return c.ServerBaseURL
}

func (c *config) GetStoragePath() string {
	return c.FileStoragePath
}

func InitConfig() Config {
	flag.Parse()
	conf := &config{
		ServerAddress:   *address,
		ServerBaseURL:   *baseURL,
		FileStoragePath: *storagepath,
	}
	err := env.Parse(conf)
	if err != nil {
		log.Fatalf("missing required env vars: %v", err)
	}
	conf.initWithFlags()
	return conf
}

func (c *config) initWithFlags() {
	if len(c.ServerAddress) == 0 {
		c.ServerAddress = *address
	}
	if len(c.ServerBaseURL) == 0 {
		c.ServerBaseURL = *baseURL
	}
	if len(c.FileStoragePath) == 0 {
		c.FileStoragePath = *storagepath
	}
}
