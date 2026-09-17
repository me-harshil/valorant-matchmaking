package rating

import (
	"context"
	"fmt"
)

type MatchResult struct {
	TeamA      []PlayerInput
	TeamB      []PlayerInput
	TeamAStats []PerformanceStats
	TeamBStats []PerformanceStats
	TeamAWon   bool
	IsDraw     bool
}

type FinalRating struct {
	PlayerID              string
	ParticipantID         string
	MuBefore, SigmaBefore float64
	MuAfter, SigmaAfter   float64
}

func ProcessMatch(m MatchResult) []FinalRating {
	rawA, rawB := UpdateRatings(m.TeamA, m.TeamB, m.TeamAWon, m.IsDraw)

	multA := NormalizeWithinTeam(m.TeamAStats)
	multB := NormalizeWithinTeam(m.TeamBStats)

	results := make([]FinalRating, 0, len(rawA)+len(rawB))
	results = append(results, applyTeam(m.TeamA, rawA, multA)...)
	results = append(results, applyTeam(m.TeamB, rawB, multB)...)
	return results
}

func applyTeam(inputs []PlayerInput, raw []PlayerOutput, mult map[string]float64) []FinalRating {
	byID := make(map[string]PlayerInput)
	for _, p := range inputs {
		byID[p.PlayerID] = p
	}

	out := make([]FinalRating, 0, len(raw))
	for _, r := range raw {
		old := byID[r.PlayerID]
		adjustedMu := ApplyPerformanceWeight(old.Mu, r.Mu, mult[r.PlayerID])
		out = append(out, FinalRating{
			PlayerID:      r.PlayerID,
			ParticipantID: old.ParticipantID,
			MuBefore:      old.Mu,
			SigmaBefore:   old.Sigma,
			MuAfter:       adjustedMu,
			SigmaAfter:    r.Sigma,
		})
	}
	return out
}

func PersistResults(ctx context.Context, results []FinalRating, history *HistoryStore, client *AccountClient) error {
	for _, r := range results {
		if err := history.Insert(ctx, r.ParticipantID, r); err != nil {
			return fmt.Errorf("failed to write rating_history for participant %s: %w", r.ParticipantID, err)
		}

		if err := client.UpdateRating(r.PlayerID, r.MuAfter, r.SigmaAfter); err != nil {
			return fmt.Errorf("rating_history recorded for %s but account-service update FAILED: %w", r.PlayerID, err)
		}
	}
	return nil
}
