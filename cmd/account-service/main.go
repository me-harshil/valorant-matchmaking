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
	"github.com/me-harshil/valorant-matchmaking/internal/account"
)

type updateRatingRequest struct {
	Mu    float64 `json:"mu"`
	Sigma float64 `json:"sigma"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()
	store, err := account.NewStore(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer store.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Put("/players/{playerID}/rating", func(w http.ResponseWriter, req *http.Request) {
		playerID := chi.URLParam(req, "playerID")

		var body updateRatingRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := store.UpdateRating(req.Context(), playerID, body.Mu, body.Sigma); err != nil {
			http.Error(w, "failed to update rating", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	r.Get("/players/{playerID}", func(w http.ResponseWriter, req *http.Request) {
		playerID := chi.URLParam(req, "playerID")

		p, err := store.GetPlayer(req.Context(), playerID)
		if err != nil {
			http.Error(w, "player not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	})

	log.Println("account-service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
