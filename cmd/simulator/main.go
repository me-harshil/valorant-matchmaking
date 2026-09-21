package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Player struct {
	PlayerID string
	Username string
	Mu       float64
	Sigma    float64
}

func (p Player) ConservativeRating() float64 {
	return p.Mu - 3*p.Sigma
}

var (
	accountURL     string
	matchmakingURL string
	matchURL       string
	ratingURL      string
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real environment variables")
	}

	accountURL = os.Getenv("ACCOUNT_SERVICE_URL")
	matchmakingURL = os.Getenv("MATCHMAKING_SERVICE_URL")
	matchURL = os.Getenv("MATCH_SERVICE_URL")
	ratingURL = os.Getenv("RATING_SERVICE_URL")

	if accountURL == "" || matchmakingURL == "" || matchURL == "" || ratingURL == "" {
		log.Fatal("all of ACCOUNT_SERVICE_URL, MATCHMAKING_SERVICE_URL, MATCH_SERVICE_URL, RATING_SERVICE_URL must be set")
	}

	const poolSize = 30 // total simulated players in the pool

	players, err := ensurePlayerPool(poolSize)
	if err != nil {
		log.Fatalf("failed to set up player pool: %v", err)
	}
	log.Printf("Player pool ready: %d players", len(players))
}

func ensurePlayerPool(poolSize int) ([]Player, error) {
	existing, err := fetchAllPlayers()
	if err != nil {
		return nil, err
	}

	switch {
	case len(existing) == poolSize:
		return existing, nil

	case len(existing) > poolSize:
		rand.Shuffle(len(existing), func(i, j int) {
			existing[i], existing[j] = existing[j], existing[i]
		})
		return existing[:poolSize], nil

	default: // len(existing) < poolSize
		needed := poolSize - len(existing)
		for i := 0; i < needed; i++ {
			username := fmt.Sprintf("sim_player_%d_%d", time.Now().UnixNano(), i)
			newPlayer, err := createPlayer(username)
			if err != nil {
				return nil, fmt.Errorf("failed to create player %s: %w", username, err)
			}
			existing = append(existing, newPlayer)
		}
		return existing, nil
	}
}

func fetchAllPlayers() ([]Player, error) {
	resp, err := http.Get(accountURL + "/players")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("account-service returned status %d", resp.StatusCode)
	}

	var raw []struct {
		PlayerID string  `json:"player_id"`
		Username string  `json:"username"`
		Mu       float64 `json:"mu"`
		Sigma    float64 `json:"sigma"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	players := make([]Player, len(raw))
	for i, r := range raw {
		players[i] = Player{PlayerID: r.PlayerID, Username: r.Username, Mu: r.Mu, Sigma: r.Sigma}
	}
	return players, nil
}

func createPlayer(username string) (Player, error) {
	body, err := json.Marshal(map[string]string{"username": username})
	if err != nil {
		return Player{}, err
	}

	resp, err := http.Post(accountURL+"/players", "application/json", bytes.NewReader(body))
	if err != nil {
		return Player{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Player{}, fmt.Errorf("account-service returned status %d", resp.StatusCode)
	}

	var result struct {
		PlayerID string `json:"player_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Player{}, err
	}

	return Player{PlayerID: result.PlayerID, Username: username, Mu: 25.0, Sigma: 8.333}, nil
}
