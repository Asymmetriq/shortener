package repository

import (
	"fmt"
	"strings"

	"github.com/Asymmetriq/shortener/internal/config"
)

type Repository interface {
	SetURL(url, userID, host string) string
	GetURL(id string) (string, error)
	GetAllURLs(userID string) ([]Data, error)
	Close() error
}

func NewRepository(cfg config.Config) Repository {
	if filepath := cfg.GetStoragePath(); len(filepath) != 0 {
		return newFileRepository(filepath, cfg.GetBaseURL())
	}
	return newInMemoryRepository()
}

type Data struct {
	userID      string
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func buildURL(host, shortURL string) string {
	if strings.Contains(host, "http") {
		return fmt.Sprintf("%s/%s", host, shortURL)
	}
	return fmt.Sprintf("http://%s/%s", host, shortURL)
}
