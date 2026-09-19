package match

import (
	"context"
	"testing"

	"github.com/me-harshil/valorant-matchmaking/internal/testutil"
)

func TestCreateMatchAndSubmitStats(t *testing.T) {
	ctx := context.Background()
	connString := testutil.SetupTestDB(t, "../../migrations/0001_init.sql")

	store, err := NewStore(ctx, connString)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	defer store.Close()

	// Insert two real players directly into the database for testing.
	var aliceID, bobID string
	err = store.pool.QueryRow(ctx, `INSERT INTO players (username) VALUES ('alice') RETURNING player_id`).Scan(&aliceID)
	if err != nil {
		t.Fatalf("insert alice failed: %v", err)
	}
	err = store.pool.QueryRow(ctx, `INSERT INTO players (username) VALUES ('bob') RETURNING player_id`).Scan(&bobID)
	if err != nil {
		t.Fatalf("insert bob failed: %v", err)
	}

	matchID, err := store.CreateMatch(ctx, CreateMatchInput{
		Map:   "Ascent",
		TeamA: []ParticipantInput{{PlayerID: aliceID, Character: "Jett"}},
		TeamB: []ParticipantInput{{PlayerID: bobID, Character: "Sova"}},
	})
	if err != nil {
		t.Fatalf("CreateMatch failed: %v", err)
	}
	if matchID == "" {
		t.Fatal("expected non-empty matchID")
	}

	err = store.SubmitStats(ctx, matchID, []SubmitStatsInput{
		{PlayerID: aliceID, Result: "win", Character: "Jett", Kills: 20, Deaths: 5, Assists: 3},
		{PlayerID: bobID, Result: "loss", Character: "Sova", Kills: 5, Deaths: 20, Assists: 1},
	})
	if err != nil {
		t.Fatalf("SubmitStats failed: %v", err)
	}

	var status string
	err = store.pool.QueryRow(ctx, `SELECT status FROM matches WHERE match_id = $1`, matchID).Scan(&status)
	if err != nil {
		t.Fatalf("failed to check match status: %v", err)
	}
	if status != "completed" {
		t.Errorf("got status %q, want completed", status)
	}

	var aliceKills int
	err = store.pool.QueryRow(ctx,
		`SELECT kills FROM match_participants WHERE match_id = $1 AND player_id = $2`,
		matchID, aliceID,
	).Scan(&aliceKills)
	if err != nil {
		t.Fatalf("failed to check alice's kills: %v", err)
	}
	if aliceKills != 20 {
		t.Errorf("got %d kills, want 20", aliceKills)
	}
}
