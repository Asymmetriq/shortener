package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Asymmetriq/shortener/internal/shorten"
)

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		storage: make(map[string]Data),
	}
}

type inMemoryRepository struct {
	storage map[string]Data
}

func (imr *inMemoryRepository) SetURL(ctx context.Context, url, userID, host string) (string, error) {
	id := shorten.Shorten(url)
	shortURL := buildURL(host, id)
	imr.storage[id] = Data{
		OriginalURL: url,
		ShortURL:    shortURL,
		userID:      userID,
	}
	return shortURL, nil
}

func (imr *inMemoryRepository) GetURL(ctx context.Context, id string) (string, error) {
	if ogURL, ok := imr.storage[id]; ok {
		return ogURL.OriginalURL, nil
	}
	return "", fmt.Errorf("no original url found with shortcut %q", id)
}

func (imr *inMemoryRepository) GetAllURLs(ctx context.Context, userID string) ([]Data, error) {
	data := make([]Data, 0)
	for _, v := range imr.storage {
		if v.userID == userID {
			data = append(data, v)
		}
	}
	if len(data) == 0 {
		return nil, errors.New("no urls for user")
	}
	return data, nil
}

func (imr *inMemoryRepository) Close() error {
	return nil
}

func (imr *inMemoryRepository) PingContext(ctx context.Context) error {
	return nil
}
