package repository

import (
	"context"
	"database/sql"

	"github.com/Asymmetriq/shortener/internal/config"
	"github.com/Asymmetriq/shortener/internal/models"
)

type Repository interface {
	SetURL(ctx context.Context, entry models.StorageEntry) error
	SetBatchURLs(ctx context.Context, entry []models.StorageEntry) error
	GetURL(ctx context.Context, id string) (string, error)
	GetAllURLs(ctx context.Context, userID string) ([]models.StorageEntry, error)
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
