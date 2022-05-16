package shortener

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Service) getHandler(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "id")

	ogURL, err := s.Storage.Get(shortID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, ogURL, http.StatusTemporaryRedirect)
}

func (s *Service) postHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(b) == 0 {
		http.Error(w, "no request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(s.getShortenedURL(r.Host, string(b))))
}

func (s *Service) jsonHandler(w http.ResponseWriter, r *http.Request) {
	var result struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(struct {
		Result string `json:"result"`
	}{
		Result: s.getShortenedURL(r.Host, result.URL),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (s *Service) getShortenedURL(host, originalURL string) string {
	shortURL := s.Storage.Set(originalURL)
	if u := s.Config.GetBaseURL(); len(u) != 0 {
		return fmt.Sprintf("%s/%s", u, shortURL)
	}
	return fmt.Sprintf("http://%s/%s", host, shortURL)
}
