package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"sync"
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

	const poolSize = 30
	players, err := ensurePlayerPool(poolSize)
	if err != nil {
		log.Fatalf("failed to set up player pool: %v", err)
	}
	log.Printf("Player pool ready: %d players", len(players))

	pool := NewPlayerPool(players)

	go runQueueingLoop(pool)
	go runMatchProcessingLoop(pool)

	log.Println("Simulator running. Press Ctrl+C to stop.")
	select {} // block forever
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

type PlayerPool struct {
	mu        sync.Mutex
	players   []Player
	available map[string]bool
}

func NewPlayerPool(players []Player) *PlayerPool {
	available := make(map[string]bool)
	for _, p := range players {
		available[p.PlayerID] = true
	}
	return &PlayerPool{players: players, available: available}
}

// TakeAvailable returns up to n available players and marks them unavailable.
func (pp *PlayerPool) TakeAvailable(n int) []Player {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	var taken []Player
	for _, p := range pp.players {
		if len(taken) >= n {
			break
		}
		if pp.available[p.PlayerID] {
			pp.available[p.PlayerID] = false
			taken = append(taken, p)
		}
	}
	return taken
}

// Release marks players as available again.
func (pp *PlayerPool) Release(players []Player) {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	for _, p := range players {
		pp.available[p.PlayerID] = true
	}
}

func queuePlayer(p Player) error {
	body, _ := json.Marshal(map[string]interface{}{
		"player_id":           p.PlayerID,
		"conservative_rating": p.ConservativeRating(),
	})

	resp, err := http.Post(matchmakingURL+"/queue", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("matchmaking-service returned status %d", resp.StatusCode)
	}
	return nil
}

func runQueueingLoop(pool *PlayerPool) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		batch := pool.TakeAvailable(10)
		if len(batch) < 10 {
			// not enough available players right now - put back what we took, try again later
			pool.Release(batch)
			continue
		}

		for _, p := range batch {
			if err := queuePlayer(p); err != nil {
				log.Printf("failed to queue player %s: %v", p.Username, err)
			}
		}
		log.Printf("Queued %d players", len(batch))
	}
}

func runMatchProcessingLoop(pool *PlayerPool) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		matches, err := fetchPendingMatches()
		if err != nil {
			log.Printf("failed to fetch pending matches: %v", err)
			continue
		}

		for _, m := range matches {
			players, err := processMatch(m.MatchID)
			if err != nil {
				log.Printf("failed to process match %s: %v", m.MatchID, err)
				continue
			}
			pool.Release(players)
			log.Printf("Processed match %s, released %d players", m.MatchID, len(players))
		}
	}
}

type MatchSummary struct {
	MatchID string `json:"match_id"`
	Status  string `json:"status"`
	Map     string `json:"map"`
}

func fetchPendingMatches() ([]MatchSummary, error) {
	resp, err := http.Get(matchURL + "/matches?status=pending")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("match-service returned status %d", resp.StatusCode)
	}

	var matches []MatchSummary
	if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
		return nil, err
	}
	return matches, nil
}

type Participant struct {
	ParticipantID string `json:"participant_id"`
	PlayerID      string `json:"player_id"`
	Team          string `json:"team"`
}

func fetchParticipants(matchID string) ([]Participant, error) {
	resp, err := http.Get(fmt.Sprintf("%s/matches/%s/participants", matchURL, matchID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("match-service returned status %d", resp.StatusCode)
	}

	var participants []Participant
	if err := json.NewDecoder(resp.Body).Decode(&participants); err != nil {
		return nil, err
	}
	return participants, nil
}

func processMatch(matchID string) ([]Player, error) {
	participants, err := fetchParticipants(matchID)
	if err != nil {
		return nil, err
	}

	teamAWon := rand.IntN(2) == 0

	var statsPayload []map[string]interface{}
	var teamA, teamB []map[string]interface{}
	var releasedPlayers []Player

	for _, p := range participants {
		var kills, deaths, assists int
		var result string

		won := (p.Team == "A" && teamAWon) || (p.Team == "B" && !teamAWon)
		if won {
			kills = 15 + rand.IntN(10)
			kills = 15 + rand.IntN(10)
			deaths = rand.IntN(8)
			result = "win"
		} else {
			kills = rand.IntN(8)
			deaths = 15 + rand.IntN(10)
			result = "loss"
		}
		assists = rand.IntN(6)

		statsPayload = append(statsPayload, map[string]interface{}{
			"PlayerID": p.PlayerID, "Result": result, "Character": "Jett",
			"Kills": kills, "Deaths": deaths, "Assists": assists,
		})

		mu, sigma, err := fetchPlayerRating(p.PlayerID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch rating for %s: %w", p.PlayerID, err)
		}

		entry := map[string]interface{}{
			"PlayerID": p.PlayerID, "ParticipantID": p.ParticipantID,
			"Mu": mu, "Sigma": sigma,
		}
		if p.Team == "A" {
			teamA = append(teamA, entry)
		} else {
			teamB = append(teamB, entry)
		}

		releasedPlayers = append(releasedPlayers, Player{PlayerID: p.PlayerID})
	}

	// submit stats
	statsBody, _ := json.Marshal(statsPayload)
	resp, err := http.Post(fmt.Sprintf("%s/matches/%s/stats", matchURL, matchID), "application/json", bytes.NewReader(statsBody))
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("submit stats returned status %d", resp.StatusCode)
	}

	// call rating-service
	ratingBody, _ := json.Marshal(map[string]interface{}{
		"TeamA": teamA, "TeamB": teamB,
		"TeamAStats": extractStats(statsPayload, "A", participants),
		"TeamBStats": extractStats(statsPayload, "B", participants),
		"TeamAWon":   teamAWon, "IsDraw": false,
	})
	resp2, err := http.Post(ratingURL+"/matches/process", "application/json", bytes.NewReader(ratingBody))
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rating-service returned status %d", resp2.StatusCode)
	}

	return releasedPlayers, nil
}

func extractStats(statsPayload []map[string]interface{}, team string, participants []Participant) []map[string]interface{} {
	teamPlayerIDs := make(map[string]bool)
	for _, p := range participants {
		if p.Team == team {
			teamPlayerIDs[p.PlayerID] = true
		}
	}

	var result []map[string]interface{}
	for _, s := range statsPayload {
		if pid, ok := s["PlayerID"].(string); ok && teamPlayerIDs[pid] {
			result = append(result, map[string]interface{}{
				"PlayerID": s["PlayerID"],
				"Kills":    s["Kills"],
				"Deaths":   s["Deaths"],
				"Assists":  s["Assists"],
			})
		}
	}
	return result
}

func fetchPlayerRating(playerID string) (mu, sigma float64, err error) {
	resp, err := http.Get(fmt.Sprintf("%s/players/%s", accountURL, playerID))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("account-service returned status %d", resp.StatusCode)
	}

	var result struct {
		Mu    float64 `json:"mu"`
		Sigma float64 `json:"sigma"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}
	return result.Mu, result.Sigma, nil
}
