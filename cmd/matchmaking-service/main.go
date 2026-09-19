package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/me-harshil/valorant-matchmaking/internal/matchmaking"
)

type joinQueueRequest struct {
	PlayerID           string  `json:"player_id"`
	ConservativeRating float64 `json:"conservative_rating"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real environment variables")
	}

	matchServiceURL := os.Getenv("MATCH_SERVICE_URL")
	if matchServiceURL == "" {
		log.Fatal("MATCH_SERVICE_URL is not set")
	}
	matchClient := matchmaking.NewMatchServiceClient(matchServiceURL)

	queue := matchmaking.NewQueue()

	go matchmaking.RunMatcherLoop(queue, 2*time.Second, func(teamA, teamB []matchmaking.QueuedPlayer) error {
		matchID, err := matchClient.CreateMatch(matchmaking.RandomMap(), teamA, teamB)
		if err != nil {
			return err
		}
		log.Printf("Match created: %s (teamA=%d players, teamB=%d players)", matchID, len(teamA), len(teamB))
		return nil
	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/queue", func(w http.ResponseWriter, req *http.Request) {
		var body joinQueueRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		queue.AddPlayer(matchmaking.QueuedPlayer{
			PlayerID:           body.PlayerID,
			ConservativeRating: body.ConservativeRating,
			QueuedAt:           time.Now(),
		})
		w.WriteHeader(http.StatusNoContent)
	})

	r.Delete("/queue/{playerID}", func(w http.ResponseWriter, req *http.Request) {
		playerID := chi.URLParam(req, "playerID")
		removed := queue.RemovePlayer(playerID)
		if !removed {
			http.Error(w, "player not in queue", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("matchmaking-service listening on :8084")
	log.Fatal(http.ListenAndServe(":8084", r))
}
