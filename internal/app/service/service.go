package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Repository interface {
	Set(url string) string
	Get(id string) (string, error)
}

func NewService(repo Repository) *Service {
	s := &Service{
		Storage: repo,
		Mux:     chi.NewMux(),
	}
	s.Use(
		middleware.Recoverer,
		middleware.RealIP,
		middleware.Logger,
	)

	s.Route("/", func(r chi.Router) {
		s.Post("/", s.postHandler)
		s.Get("/{id}", s.getHandler)

		s.Post("/api/shorten", s.jsonHandler)

		// Так не работает? Почему?
		// s.Route("/api", func(r chi.Router) {
		// 	s.Post("/shorten", s.jsonHandler)
		// })
	})

	return s
}

type Service struct {
	*chi.Mux
	Storage Repository
}
