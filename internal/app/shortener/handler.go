package shortener

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Asymmetriq/shortener/internal/cookie"
	"github.com/go-chi/chi/v5"
)

func (s *Service) getHandler(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "id")

	ogURL, err := s.Storage.GetURL(shortID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, ogURL, http.StatusTemporaryRedirect)
}

func (s *Service) postHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(cookie.Name).(string)
	if !ok {
		http.Error(w, "no userID provided", http.StatusBadRequest)
		return
	}

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
	host := r.Host
	if u := s.Config.GetBaseURL(); len(u) != 0 {
		host = u
	}

	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(s.Storage.SetURL(string(b), userID, host)))
}

func (s *Service) jsonHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(cookie.Name).(string)
	if !ok {
		http.Error(w, "no userID provided", http.StatusBadRequest)
		return
	}

	var result struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	host := r.Host
	if u := s.Config.GetBaseURL(); len(u) != 0 {
		host = u
	}

	resp, err := json.Marshal(struct {
		Result string `json:"result"`
	}{
		Result: s.Storage.SetURL(result.URL, userID, host),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (s *Service) pingHandler(w http.ResponseWriter, r *http.Request) {
	err := s.DB.PingContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Service) userURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(cookie.Name).(string)
	if !ok {
		http.Error(w, "no userID provided", http.StatusBadRequest)
		return
	}

	urls, err := s.Storage.GetAllURLs(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(urls); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
