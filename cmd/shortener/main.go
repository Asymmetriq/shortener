package main

import (
	"net/http"

	"github.com/Asymmetriq/shortener/internal/app/shortener"
	"github.com/Asymmetriq/shortener/internal/config"
	"github.com/Asymmetriq/shortener/internal/repositories"
)

func main() {
	cfg := config.InitConfig()
	repo := repositories.NewRepository(cfg.GetStoragePath())
	defer repo.Close()

	service := shortener.NewShortener(repo, cfg)
	http.ListenAndServe(service.Config.GetAddress(), service)
}
