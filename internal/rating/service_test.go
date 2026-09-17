package rating

import "testing"

func TestProcessMatch_FiveOnFive(t *testing.T) {
	teamA := []PlayerInput{
		{"a1", 25, 8.333}, {"a2", 25, 8.333}, {"a3", 25, 8.333},
		{"a4", 25, 8.333}, {"a5", 25, 8.333},
	}
	teamB := []PlayerInput{
		{"b1", 25, 8.333}, {"b2", 25, 8.333}, {"b3", 25, 8.333},
		{"b4", 25, 8.333}, {"b5", 25, 8.333},
	}

	// a1 is the standout performer on team A
	statsA := []PerformanceStats{
		{"a1", 25, 5, 8}, {"a2", 10, 8, 3}, {"a3", 9, 9, 2},
		{"a4", 8, 10, 4}, {"a5", 7, 11, 1},
	}
	statsB := []PerformanceStats{
		{"b1", 8, 12, 2}, {"b2", 7, 13, 1}, {"b3", 9, 11, 3},
		{"b4", 6, 14, 0}, {"b5", 5, 15, 2},
	}

	results := ProcessMatch(MatchResult{
		TeamA: teamA, TeamB: teamB,
		TeamAStats: statsA, TeamBStats: statsB,
		TeamAWon: true, IsDraw: false,
	})

	if len(results) != 10 {
		t.Fatalf("got %d results, want 10", len(results))
	}

	var a1, a5 FinalRating
	for _, r := range results {
		if r.PlayerID == "a1" {
			a1 = r
		}
		if r.PlayerID == "a5" {
			a5 = r
		}
	}

	// a1 (standout on winning team) should gain MORE mu than a5 (weakest on winning team)
	a1Gain := a1.MuAfter - a1.MuBefore
	a5Gain := a5.MuAfter - a5.MuBefore
	if a1Gain <= a5Gain {
		t.Errorf("expected standout player (a1 gain=%v) to out-gain weak player (a5 gain=%v)", a1Gain, a5Gain)
	}
}
