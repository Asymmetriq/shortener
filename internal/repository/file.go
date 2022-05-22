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
	OriginalURL string
	ShortURL    string
}

type fileRepostitory struct {
	file    *os.File
	encoder *json.Encoder
	storage map[string]string
}

func newFileRepository(filename string) *fileRepostitory {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("create file repo: %v", err)
	}
	data, err := restoreData(f)
	if err != nil {
		log.Fatalf("corrupted data %v", err)
	}

	return &fileRepostitory{
		file:    f,
		encoder: json.NewEncoder(f),
		storage: data,
	}
}

func (fr *fileRepostitory) Set(url string) string {
	shortURL := shorten.Shorten(url)
	if _, ok := fr.storage[shortURL]; !ok {
		fr.encoder.Encode(dataJSON{
			OriginalURL: url,
			ShortURL:    shortURL,
		})
	}
	fr.storage[shortURL] = url

	return shortURL
}

func (fr *fileRepostitory) Get(id string) (string, error) {
	if ogURL, ok := fr.storage[id]; ok {
		return ogURL, nil
	}
	return "", fmt.Errorf("no original url found with shortcut %q", id)
}

func (fr *fileRepostitory) Close() error {
	return fr.file.Close()
}

func restoreData(file *os.File) (map[string]string, error) {
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stats.Size() == 0 {
		return make(map[string]string), nil
	}

	restored := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		v := dataJSON{}
		if err = json.Unmarshal(scanner.Bytes(), &v); err != nil {
			return nil, err
		}
		restored[v.ShortURL] = v.OriginalURL
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return restored, nil
}
