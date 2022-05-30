package shortener

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Asymmetriq/shortener/internal/cookie"
	"github.com/Asymmetriq/shortener/internal/models"
	"github.com/go-chi/chi/v5"
)

func (s *Service) getHandler(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "id")

	ogURL, err := s.Storage.GetURL(r.Context(), shortID)
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

	entry := models.NewStorageEntry(string(b), host, userID)
	err = s.Storage.SetURL(r.Context(), entry)
	code := models.ParseStorageError(err)
	if code == http.StatusBadRequest {
		http.Error(w, err.Error(), code)
		return
	}
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(code)
	w.Write([]byte(entry.ShortURL))
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
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	host := r.Host
	if u := s.Config.GetBaseURL(); len(u) != 0 {
		host = u
	}

	entry := models.NewStorageEntry(result.URL, host, userID)
	err = s.Storage.SetURL(r.Context(), entry)
	code := models.ParseStorageError(err)
	if code == http.StatusBadRequest {
		http.Error(w, err.Error(), code)
		return
	}
	resp, err := json.Marshal(struct {
		Result string `json:"result"`
	}{
		Result: entry.ShortURL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func (s *Service) batchHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(cookie.Name).(string)
	if !ok {
		http.Error(w, "no userID provided", http.StatusBadRequest)
		return
	}

	var entries []models.StorageEntry
	err := json.NewDecoder(r.Body).Decode(&entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	host := r.Host
	if u := s.Config.GetBaseURL(); len(u) != 0 {
		host = u
	}

	for i := range entries {
		entries[i].UserID = userID
		if err = entries[i].BuildShortURL(host); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = s.Storage.SetBatchURLs(r.Context(), entries)
	code := models.ParseStorageError(err)
	if code == http.StatusBadRequest {
		http.Error(w, err.Error(), code)
		return
	}
	for i := range entries {
		entries[i].ID = ""
		entries[i].OriginalURL = ""
		entries[i].UserID = ""
	}

	value, err := json.Marshal(entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(value)
}

func (s *Service) userURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(cookie.Name).(string)
	if !ok {
		http.Error(w, "no userID provided", http.StatusBadRequest)
		return
	}

	urls, err := s.Storage.GetAllURLs(r.Context(), userID)
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

func (s *Service) pingHandler(w http.ResponseWriter, r *http.Request) {
	err := s.Storage.PingContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
