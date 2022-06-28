package main

import (
	"net/http"

	"github.com/Asymmetriq/shortener/internal/app/shortener"
	"github.com/Asymmetriq/shortener/internal/config"
	"github.com/Asymmetriq/shortener/internal/database"
	"github.com/Asymmetriq/shortener/internal/repository"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	cfg := config.InitConfig()
	db := database.ConnectToDatabase("pgx", cfg.GetDatabaseDSN())

	repo := repository.NewRepository(cfg, db)
	defer repo.Close()

	service := shortener.NewShortener(repo, cfg)
	http.ListenAndServe(service.Config.GetAddress(), service)
}
