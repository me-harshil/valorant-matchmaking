package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/me-harshil/valorant-matchmaking/internal/rating"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	history := rating.NewHistoryStore(pool)

	accountServiceURL := os.Getenv("ACCOUNT_SERVICE_URL")
	if accountServiceURL == "" {
		log.Fatal("ACCOUNT_SERVICE_URL is not set")
	}
	client := rating.NewAccountClient(accountServiceURL)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/ratings/process", func(w http.ResponseWriter, req *http.Request) {
		var m rating.MatchResult
		if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		results := rating.ProcessMatch(m)

		if err := rating.PersistResults(req.Context(), results, history, client); err != nil {
			http.Error(w, "failed to persist ratings: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	log.Println("rating-service listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
