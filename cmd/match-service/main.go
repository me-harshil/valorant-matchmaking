package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/me-harshil/valorant-matchmaking/internal/match"
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
	store, err := match.NewStore(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer store.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/matches", func(w http.ResponseWriter, req *http.Request) {
		var input match.CreateMatchInput
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		matchID, err := store.CreateMatch(req.Context(), input)
		if err != nil {
			http.Error(w, "failed to create match: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"match_id": matchID})
	})

	r.Post("/matches/{matchID}/stats", func(w http.ResponseWriter, req *http.Request) {
		matchID := chi.URLParam(req, "matchID")

		var stats []match.SubmitStatsInput
		if err := json.NewDecoder(req.Body).Decode(&stats); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := store.SubmitStats(req.Context(), matchID, stats); err != nil {
			http.Error(w, "failed to submit stats: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("match-service listening on :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}
