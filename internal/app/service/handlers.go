package service

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
	shortURL := s.Storage.Set(string(b))

	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://%s/%s", r.Host, shortURL)))
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

	shortURL := s.Storage.Set(string(result.URL))
	resp, err := json.Marshal(struct {
		Result string `json:"result"`
	}{Result: fmt.Sprintf("http://%s/%s", r.Host, shortURL)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
