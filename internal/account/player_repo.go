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

type PlayerSummary struct {
	PlayerID string  `json:"player_id"`
	Username string  `json:"username"`
	Mu       float64 `json:"mu"`
	Sigma    float64 `json:"sigma"`
}

func (s *Store) ListPlayers(ctx context.Context) ([]PlayerSummary, error) {
	rows, err := s.pool.Query(ctx, `SELECT player_id, username, mu, sigma FROM players`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PlayerSummary
	for rows.Next() {
		var p PlayerSummary
		if err := rows.Scan(&p.PlayerID, &p.Username, &p.Mu, &p.Sigma); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, rows.Err()
}
