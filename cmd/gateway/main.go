package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func newProxy(target string) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatalf("invalid target URL %s: %v", target, err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real environment variables")
	}

	accountURL := os.Getenv("ACCOUNT_SERVICE_URL")
	matchURL := os.Getenv("MATCH_SERVICE_URL")
	ratingURL := os.Getenv("RATING_SERVICE_URL")
	matchmakingURL := os.Getenv("MATCHMAKING_SERVICE_URL")

	if accountURL == "" || matchURL == "" || ratingURL == "" || matchmakingURL == "" {
		log.Fatal("all 4 service URL env vars must be set")
	}

	accountProxy := newProxy(accountURL)
	matchProxy := newProxy(matchURL)
	ratingProxy := newProxy(ratingURL)
	matchmakingProxy := newProxy(matchmakingURL)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	r.Handle("/players*", accountProxy)
	r.Handle("/matches*", matchProxy)
	r.Handle("/queue*", matchmakingProxy)
	r.Handle("/ratings*", ratingProxy)

	log.Println("gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
