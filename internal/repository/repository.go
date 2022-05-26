package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Asymmetriq/shortener/internal/config"
)

type Repository interface {
	SetURL(ctx context.Context, url, userID, host string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
	GetAllURLs(ctx context.Context, userID string) ([]Data, error)
	Close() error
	PingContext(ctx context.Context) error
}

func NewRepository(cfg config.Config, db *sql.DB) Repository {
	if db != nil {
		return newDBRepository(db)
	}
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
