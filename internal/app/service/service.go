package service

import (
	"github.com/Asymmetriq/shortener/internal/app/repository"
)

type Repository interface {
	Set(url []byte) string
	Get(id string) (string, error)
}

func NewService() *Service {
	return &Service{
		Storage: repository.NewRepository(),
	}
}

type Service struct {
	Storage Repository
}
