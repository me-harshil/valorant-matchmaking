package rating

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/me-harshil/valorant-matchmaking/internal/testutil"
)

func TestTierFromRating(t *testing.T) {
	ctx := context.Background()
	connString := testutil.SetupTestDB(t,
		"../../migrations/0001_init.sql",
		"../../migrations/0002_rank_tiers.sql",
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer pool.Close()

	ts := NewTierStore(pool)

	tier, err := ts.TierFromRating(ctx, 47.0)
	if err != nil {
		t.Fatalf("TierFromRating failed: %v", err)
	}
	if tier.Name != "Gold" || tier.Sub != 1 {
		t.Errorf("got %s %d, want Gold 1", tier.Name, tier.Sub)
	}
}
