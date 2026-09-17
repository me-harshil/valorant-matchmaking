package rating

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Tier struct {
	Name string
	Sub  int
}

type TierStore struct {
	pool *pgxpool.Pool
}

func NewTierStore(pool *pgxpool.Pool) *TierStore {
	return &TierStore{pool: pool}
}

// TierFromRating finds the highest tier whose min_rating is <= conservativeRating.
func (ts *TierStore) TierFromRating(ctx context.Context, conservativeRating float64) (Tier, error) {
	var t Tier
	err := ts.pool.QueryRow(ctx,
		`SELECT name, sub FROM rank_tiers
		 WHERE min_rating <= $1
		 ORDER BY min_rating DESC
		 LIMIT 1`,
		conservativeRating,
	).Scan(&t.Name, &t.Sub)
	return t, err
}
