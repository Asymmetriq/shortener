package main

import (
	"net/http"

	repo "github.com/Asymmetriq/shortener/internal/app/repository"
	"github.com/Asymmetriq/shortener/internal/app/service"
)

func main() {

	service := service.NewService(repo.NewRepository())
	http.ListenAndServe(":8080", service)
}
