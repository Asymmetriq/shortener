package repository

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
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
	h := sha1.New()
	h.Write(url)
	shortURL := hex.EncodeToString(h.Sum(nil))
	imr.storage[shortURL] = string(url)
	return shortURL
}

func (imr *inMemoryRepository) Get(id string) (string, error) {
	if ogURL, ok := imr.storage[id]; ok {
		return ogURL, nil
	}
	return "", fmt.Errorf("no original url found with shortcut %q", id)
}
