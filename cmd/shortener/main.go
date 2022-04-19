package main

import (
	"net/http"

	"github.com/Asymmetriq/shortener/internal/app/service"
)

func main() {
	service := service.NewService()

	http.HandleFunc("/", service.Multiplexer)

	http.ListenAndServe(":8080", nil)
}
