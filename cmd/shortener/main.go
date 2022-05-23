package main

import (
	"net/http"

	"github.com/Asymmetriq/shortener/internal/app/shortener"
	"github.com/Asymmetriq/shortener/internal/config"
	"github.com/Asymmetriq/shortener/internal/db"
	"github.com/Asymmetriq/shortener/internal/repository"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	cfg := config.InitConfig()

	repo := repository.NewRepository(cfg)
	defer repo.Close()

	con := db.ConnectToPostgres("pgx", cfg.GetDatabaseDSN())
	defer con.Close()

	service := shortener.NewShortener(repo, cfg, con)
	http.ListenAndServe(service.Config.GetAddress(), service)
}
