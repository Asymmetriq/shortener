package main

import (
	"net/http"

	"github.com/Asymmetriq/shortener/internal/app/shortener"
	"github.com/Asymmetriq/shortener/internal/config"
	"github.com/Asymmetriq/shortener/internal/repository"
)

func main() {
	cfg := config.InitConfig()

	repo := repository.NewRepository(cfg)
	defer repo.Close()

	service := shortener.NewShortener(repo, cfg)
	http.ListenAndServe(service.Config.GetAddress(), service)
}
