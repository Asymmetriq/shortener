package repository

import (
	"context"

	"github.com/Asymmetriq/shortener/internal/config"
	"github.com/Asymmetriq/shortener/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	SetURL(ctx context.Context, entry models.StorageEntry) error
	SetBatchURLs(ctx context.Context, entry []models.StorageEntry) error
	GetURL(ctx context.Context, id string) (string, error)
	GetAllURLs(ctx context.Context, userID string) ([]models.StorageEntry, error)
	BatchDelete(ctx context.Context, req models.DeleteRequest)
	Close() error
	PingContext(ctx context.Context) error
}

func NewRepository(cfg config.Config, db *sqlx.DB) Repository {
	if db != nil {
		return newDBRepository(db)
	}
	if filepath := cfg.GetStoragePath(); len(filepath) != 0 {
		return newFileRepository(filepath, cfg.GetBaseURL())
	}
	return newInMemoryRepository()
}
