package shortener

import (
	"github.com/Asymmetriq/shortener/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Repository interface {
	Set(url string) string
	Get(id string) (string, error)
}

func NewShortener(repo Repository, cfg config.Config) *Service {
	s := &Service{
		Mux:     chi.NewMux(),
		Storage: repo,
		Config:  cfg,
	}

	s.Use(
		middleware.Recoverer,
		middleware.RealIP,
		middleware.Logger,
	)
	s.Route("/", func(r chi.Router) {
		s.Post("/", s.postHandler)
		s.Get("/{id}", s.getHandler)

		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", s.jsonHandler)
		})
	})

	return s
}

type Service struct {
	*chi.Mux
	Storage Repository
	Config  config.Config
}
