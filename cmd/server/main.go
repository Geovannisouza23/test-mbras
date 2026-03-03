package main

import (
	"log"
	"net/http"

	"backend-challenge-092025/internal/app"
	"backend-challenge-092025/internal/domain"
)

func main() {
	processor := domain.NewProcessor()
	handler := app.NewHandler(processor)
	http.HandleFunc("/analyze-feed", handler.AnalyzeFeed)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
