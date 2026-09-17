package rating

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HistoryStore struct {
	pool *pgxpool.Pool
}

func NewHistoryStore(pool *pgxpool.Pool) *HistoryStore {
	return &HistoryStore{pool: pool}
}

func (h *HistoryStore) Insert(ctx context.Context, participantID string, r FinalRating) error {
	_, err := h.pool.Exec(ctx,
		`INSERT INTO rating_history (participant_id, mu_before, sigma_before, mu_after, sigma_after)
		 VALUES ($1, $2, $3, $4, $5)`,
		participantID, r.MuBefore, r.SigmaBefore, r.MuAfter, r.SigmaAfter,
	)
	return err
}
