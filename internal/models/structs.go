package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Asymmetriq/shortener/internal/shorten"
)

type StorageEntry struct {
	ID            string `json:"id,omitempty" db:"id"`
	UserID        string `json:"user_id,omitempty" db:"user_id"`
	ShortURL      string `json:"short_url" db:"short_url"`
	OriginalURL   string `json:"original_url,omitempty" db:"original_url"`
	CorrelationID string `json:"correlation_id,omitempty"`
	Deleted       bool   `json:"-" db:"deleted"`
}

func NewStorageEntry(originalURL, host, userID string) StorageEntry {
	id := shorten.Shorten(originalURL)
	shortURL := buildURL(host, id)

	return StorageEntry{
		ID:          id,
		UserID:      userID,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}

}

func (e *StorageEntry) BuildShortURL(host string) error {
	if len(e.OriginalURL) == 0 {
		return errors.New("can't build url, original url empty")
	}
	e.ID = shorten.Shorten(e.OriginalURL)
	e.ShortURL = buildURL(host, e.ID)
	return nil
}

func buildURL(host, id string) string {
	if strings.Contains(host, "http") {
		return fmt.Sprintf("%s/%s", host, id)
	}
	return fmt.Sprintf("http://%s/%s", host, id)
}

type DeleteRequest struct {
	UserID string
	IDs    []string
}
