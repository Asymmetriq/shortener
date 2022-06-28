package shortener

import (
	"github.com/Asymmetriq/shortener/internal/config"
	repo "github.com/Asymmetriq/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewShortener(repo repo.Repository, cfg config.Config) *Service {
	s := &Service{
		Mux:     chi.NewMux(),
		Storage: repo,
		Config:  cfg,
	}

	s.Use(
		middleware.Recoverer,
		middleware.RealIP,
		middleware.Logger,

		gzipMiddleware,
		cookieMiddleware,
	)
	s.Route("/", func(r chi.Router) {
		s.Post("/", s.postHandler)
		s.Get("/{id}", s.getHandler)

		s.Get("/ping", s.pingHandler)

		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", s.jsonHandler)
				r.Post("/batch", s.batchHandler)
			})

			r.Route("/user", func(r chi.Router) {
				r.Get("/urls", s.userURLsHandler)
				r.Delete("/urls", s.asyncDeleteHandler)
			})
		})
	})

	return s
}

type Service struct {
	*chi.Mux
	Storage repo.Repository
	Config  config.Config
}
