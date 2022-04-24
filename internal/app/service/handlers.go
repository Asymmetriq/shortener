package service

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Service) Multiplexer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.getHandler(w, r)
		return
	}
	s.postHandler(w, r)
}

func (s *Service) getHandler(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "id")

	if len(r.URL.Path) == 0 {
		http.Error(w, "empty url", 400)
		return
	}
	ogURL, err := s.Storage.Get(shortID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	http.Redirect(w, r, ogURL, http.StatusTemporaryRedirect)
}

func (s *Service) postHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if len(b) == 0 {
		http.Error(w, "no request body", 400)
		return
	}
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://%s/%s", r.Host, s.Storage.Set(b))))

}
