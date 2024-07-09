package main

import (
	"log"
	"net/http"

	handler "lotteryapi/api"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/get-results", handler.GetResults)
	r.Get("/check-results", handler.CheckResults)
	r.Get("/analyze-results", handler.AnalyzeResults)
	r.Get("/get-lotteries", handler.GetLotteries)
	r.Get("/hello", handler.Hello)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
