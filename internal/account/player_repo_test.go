package account

import (
	"context"
	"testing"

	"github.com/me-harshil/valorant-matchmaking/internal/testutil"
)

func TestUpdateAndGetRating(t *testing.T) {
	ctx := context.Background()
	connString := testutil.SetupTestDB(t, "../../migrations/0001_init.sql")

	store, err := NewStore(ctx, connString)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	defer store.Close()

	var playerID string
	err = store.pool.QueryRow(ctx,
		`INSERT INTO players (username) VALUES ($1) RETURNING player_id`,
		"test_user",
	).Scan(&playerID)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	err = store.UpdateRating(ctx, playerID, 30.5, 6.2)
	if err != nil {
		t.Fatalf("UpdateRating failed: %v", err)
	}

	got, err := store.GetPlayer(ctx, playerID)
	if err != nil {
		t.Fatalf("GetPlayer failed: %v", err)
	}

	if got.Mu != 30.5 || got.Sigma != 6.2 {
		t.Errorf("got mu=%v sigma=%v, want mu=30.5 sigma=6.2", got.Mu, got.Sigma)
	}
}
