package matchmaking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type MatchServiceClient struct {
	baseURL string
	client  *http.Client
}

func NewMatchServiceClient(baseURL string) *MatchServiceClient {
	return &MatchServiceClient{baseURL: baseURL, client: &http.Client{}}
}

type createMatchParticipant struct {
	PlayerID  string
	Character string
}

type createMatchBody struct {
	Map   string
	TeamA []createMatchParticipant
	TeamB []createMatchParticipant
}

func (c *MatchServiceClient) CreateMatch(mapName string, teamA, teamB []QueuedPlayer) (string, error) {
	toParticipants := func(team []QueuedPlayer) []createMatchParticipant {
		out := make([]createMatchParticipant, len(team))
		for i, p := range team {
			out[i] = createMatchParticipant{PlayerID: p.PlayerID}
		}
		return out
	}

	body, err := json.Marshal(createMatchBody{
		Map:   mapName,
		TeamA: toParticipants(teamA),
		TeamB: toParticipants(teamB),
	})
	if err != nil {
		return "", err
	}

	resp, err := c.client.Post(c.baseURL+"/matches", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("match-service returned status %d", resp.StatusCode)
	}

	var result struct {
		MatchID string `json:"match_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.MatchID, nil
}
