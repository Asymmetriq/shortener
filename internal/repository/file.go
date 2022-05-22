package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Asymmetriq/shortener/internal/shorten"
)

type dataJSON struct {
	OriginalURL string `json:"original_url"`
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
}

type fileRepostitory struct {
	file    *os.File
	encoder *json.Encoder
	storage map[string]Data
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

func (fr *fileRepostitory) SetURL(url, userID, host string) string {
	id := shorten.Shorten(url)
	shortURL := buildURL(host, id)

	if _, ok := fr.storage[id]; !ok {
		fr.encoder.Encode(dataJSON{
			OriginalURL: url,
			ID:          id,
			UserID:      userID,
		})
	}
	fr.storage[id] = Data{
		OriginalURL: url,
		ShortURL:    shortURL,
		userID:      userID,
	}

	return shortURL
}

func (fr *fileRepostitory) GetURL(id string) (string, error) {
	if ogURL, ok := fr.storage[id]; ok {
		return ogURL.OriginalURL, nil
	}
	return "", fmt.Errorf("no original url found with shortcut %q", id)
}

func (fr *fileRepostitory) GetAllURLs(userID string) ([]Data, error) {
	data := make([]Data, 0)
	for _, v := range fr.storage {
		if v.userID == userID {
			data = append(data, v)
		}
	}
	return data, nil
}

func (fr *fileRepostitory) Close() error {
	return fr.file.Close()
}

func restoreData(file *os.File, host string) (map[string]Data, error) {
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stats.Size() == 0 {
		return make(map[string]Data), nil
	}

	restored := make(map[string]Data)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		v := dataJSON{}
		if err = json.Unmarshal(scanner.Bytes(), &v); err != nil {
			return nil, err
		}
		restored[v.ID] = Data{
			OriginalURL: v.OriginalURL,
			ShortURL:    buildURL(host, v.ID),
			userID:      v.UserID,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return restored, nil
}
