package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Asymmetriq/shortener/internal/models"
)

type fileRepostitory struct {
	file    *os.File
	encoder *json.Encoder
	storage map[string]models.StorageEntry
}

func newFileRepository(filename, host string) *fileRepostitory {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("create file repo: %v", err)
	}
	data, err := restoreData(f, host)
	if err != nil {
		log.Fatalf("corrupted data %v", err)
	}

	return &fileRepostitory{
		file:    f,
		encoder: json.NewEncoder(f),
		storage: data,
	}
}

func (fr *fileRepostitory) SetURL(ctx context.Context, entry models.StorageEntry) error {
	if _, ok := fr.storage[entry.ID]; !ok {
		fr.encoder.Encode(entry)
	}
	fr.storage[entry.ID] = entry
	return nil
}

func (fr *fileRepostitory) SetBatchURLs(ctx context.Context, entries []models.StorageEntry) error {
	for _, entry := range entries {
		fr.storage[entry.ID] = entry
	}
	return nil
}

func (fr *fileRepostitory) GetURL(ctx context.Context, id string) (string, error) {
	item, ok := fr.storage[id]
	if !ok {
		return "", fmt.Errorf("no original url found with shortcut %q", id)
	}
	if item.Deleted {
		return "", models.ErrDeleted
	}
	return item.OriginalURL, nil
}

func (fr *fileRepostitory) GetAllURLs(ctx context.Context, userID string) ([]models.StorageEntry, error) {
	data := make([]models.StorageEntry, 0)
	for _, v := range fr.storage {
		if v.UserID == userID {
			data = append(data, v)
		}
	}
	return data, nil
}

func (fr *fileRepostitory) BatchDelete(ctx context.Context, req models.DeleteRequest) {
	for _, v := range req.IDs {
		if item := fr.storage[v]; item.UserID == req.UserID {
			item.Deleted = true
			fr.storage[v] = item
		}
	}
}

func (fr *fileRepostitory) Close() error {
	return fr.file.Close()
}

func (fr *fileRepostitory) PingContext(ctx context.Context) error {
	return nil
}

func restoreData(file *os.File, host string) (map[string]models.StorageEntry, error) {
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}
	restored := make(map[string]models.StorageEntry)
	if stats.Size() == 0 {
		return restored, nil
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := models.StorageEntry{}
		if err = json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, err
		}
		restored[entry.ID] = entry
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return restored, nil
}
