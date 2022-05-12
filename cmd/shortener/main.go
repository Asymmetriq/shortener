package main

import (
	"net/http"

	repo "github.com/Asymmetriq/shortener/internal/app/repository"
	"github.com/Asymmetriq/shortener/internal/app/shortener"
	"github.com/Asymmetriq/shortener/internal/config"
)

func main() {
	service := shortener.NewShortener(
		repo.NewRepository(),
		config.InitConfig(),
	)
	http.ListenAndServe(service.Config.GetAddress(), service)
}
