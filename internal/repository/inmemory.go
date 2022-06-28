package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Asymmetriq/shortener/internal/models"
)

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		storage: make(map[string]models.StorageEntry),
	}
}

type inMemoryRepository struct {
	storage map[string]models.StorageEntry
}

func (imr *inMemoryRepository) SetURL(ctx context.Context, entry models.StorageEntry) error {
	imr.storage[entry.ID] = entry
	return nil
}

func (imr *inMemoryRepository) SetBatchURLs(ctx context.Context, entries []models.StorageEntry) error {
	for _, entry := range entries {
		imr.storage[entry.ID] = entry
	}
	return nil
}

func (imr *inMemoryRepository) GetURL(ctx context.Context, id string) (string, error) {
	item, ok := imr.storage[id]
	if !ok {
		return "", fmt.Errorf("no original url found with shortcut %q", id)
	}
	if item.Deleted {
		return "", models.ErrDeleted
	}
	return item.OriginalURL, nil
}

func (imr *inMemoryRepository) GetAllURLs(ctx context.Context, userID string) ([]models.StorageEntry, error) {
	data := make([]models.StorageEntry, 0)
	for _, entry := range imr.storage {
		if entry.UserID == userID {
			data = append(data, entry)
		}
	}
	if len(data) == 0 {
		return nil, errors.New("no urls for user")
	}
	return data, nil
}

func (imr *inMemoryRepository) BatchDelete(req models.DeleteRequest) {
	for _, v := range req.IDs {
		if item := imr.storage[v]; item.UserID == req.UserID {
			item.Deleted = true
			imr.storage[v] = item
		}
	}
}

func (imr *inMemoryRepository) Close() error {
	return nil
}

func (imr *inMemoryRepository) PingContext(ctx context.Context) error {
	return nil
}
