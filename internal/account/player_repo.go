package account

import "context"

type PlayerRating struct {
	PlayerID string
	Mu       float64
	Sigma    float64
}

func (s *Store) UpdateRating(ctx context.Context, playerID string, mu, sigma float64) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE players SET mu = $1, sigma = $2, updated_at = now() WHERE player_id = $3`,
		mu, sigma, playerID,
	)
	return err
}

func (s *Store) GetPlayer(ctx context.Context, playerID string) (PlayerRating, error) {
	var p PlayerRating
	err := s.pool.QueryRow(ctx,
		`SELECT player_id, mu, sigma FROM players WHERE player_id = $1`,
		playerID,
	).Scan(&p.PlayerID, &p.Mu, &p.Sigma)
	return p, err
}

func (s *Store) CreatePlayer(ctx context.Context, username string) (string, error) {
	var playerID string
	err := s.pool.QueryRow(ctx,
		`INSERT INTO players (username) VALUES ($1) RETURNING player_id`,
		username,
	).Scan(&playerID)
	return playerID, err
}
