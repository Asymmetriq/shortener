package repository

import (
	"fmt"

	"github.com/Asymmetriq/shortener/internal/app/shorten"
)

func NewRepository() *inMemoryRepository {
	return &inMemoryRepository{
		storage: make(map[string]string),
	}
}

type inMemoryRepository struct {
	storage map[string]string
}

func (imr *inMemoryRepository) Set(url []byte) string {
	shortURL := shorten.Shorten(url)
	imr.storage[shortURL] = string(url)
	return shortURL
}

func (imr *inMemoryRepository) Get(id string) (string, error) {
	if ogURL, ok := imr.storage[id]; ok {
		return ogURL, nil
	}
	return "", fmt.Errorf("no original url found with shortcut %q", id)
}
