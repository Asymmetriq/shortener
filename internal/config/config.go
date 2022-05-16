package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

func init() {
	flag.StringVar(&_address, "a", ":8080", "Server's host address")
	flag.StringVar(&_baseURL, "b", "", "Server's base URL")
	flag.StringVar(&_storagepath, "f", "", "Storage path")
}

var (
	_address     string
	_baseURL     string
	_storagepath string
)

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
		ServerAddress:   _address,
		ServerBaseURL:   _baseURL,
		FileStoragePath: _storagepath,
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
		c.ServerAddress = _address
	}
	if len(c.ServerBaseURL) == 0 {
		c.ServerBaseURL = _baseURL
	}
	if len(c.FileStoragePath) == 0 {
		c.FileStoragePath = _storagepath
	}
}
