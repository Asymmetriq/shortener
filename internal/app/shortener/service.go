package shortener

import (
	"github.com/Asymmetriq/shortener/internal/config"
	r "github.com/Asymmetriq/shortener/internal/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewShortener(repo r.Repository, cfg config.Config) *Service {
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
	Storage r.Repository
	Config  config.Config
}
