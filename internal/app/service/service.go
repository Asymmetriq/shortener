package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Repository interface {
	Set(url []byte) string
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
	})
	return s

}

type Service struct {
	*chi.Mux
	Storage Repository
}
