package service

import (
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
