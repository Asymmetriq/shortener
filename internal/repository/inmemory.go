package repository

import (
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

func (imr *inMemoryRepository) SetURL(url, userID, host string) string {
	id := shorten.Shorten(url)
	shortURL := buildURL(host, id)
	imr.storage[id] = Data{
		OriginalURL: url,
		ShortURL:    shortURL,
		userID:      userID,
	}
	return shortURL
}

func (imr *inMemoryRepository) GetURL(id string) (string, error) {
	if ogURL, ok := imr.storage[id]; ok {
		return ogURL.OriginalURL, nil
	}
	return "", fmt.Errorf("no original url found with shortcut %q", id)
}

func (imr *inMemoryRepository) GetAllURLs(userID string) ([]Data, error) {
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
