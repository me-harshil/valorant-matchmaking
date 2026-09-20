package match

import "context"

const unselectedCharacter = "unselected"

type CreateMatchInput struct {
	Map   string
	TeamA []ParticipantInput
	TeamB []ParticipantInput
}

type ParticipantInput struct {
	PlayerID  string
	Character string
}

func (s *Store) CreateMatch(ctx context.Context, input CreateMatchInput) (string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var matchID string
	err = tx.QueryRow(ctx,
		`INSERT INTO matches (status, map) VALUES ('pending', $1) RETURNING match_id`,
		input.Map,
	).Scan(&matchID)
	if err != nil {
		return "", err
	}

	for _, p := range input.TeamA {
		character := p.Character
		if character == "" {
			character = unselectedCharacter
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO match_participants (match_id, player_id, team, result, character)
			 VALUES ($1, $2, 'A', 'draw', $3)`,
			matchID, p.PlayerID, character,
		)
		if err != nil {
			return "", err
		}
	}
	for _, p := range input.TeamB {
		character := p.Character
		if character == "" {
			character = unselectedCharacter
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO match_participants (match_id, player_id, team, result, character)
			 VALUES ($1, $2, 'B', 'draw', $3)`,
			matchID, p.PlayerID, character,
		)
		if err != nil {
			return "", err
		}
	}

	return matchID, tx.Commit(ctx)
}

type SubmitStatsInput struct {
	PlayerID  string
	Result    string
	Character string
	Kills     int
	Deaths    int
	Assists   int
}

func (s *Store) SubmitStats(ctx context.Context, matchID string, stats []SubmitStatsInput) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, st := range stats {
		_, err = tx.Exec(ctx,
			`UPDATE match_participants
	 		SET result = $1, character = $2, kills = $3, deaths = $4, assists = $5
	 		WHERE match_id = $6 AND player_id = $7`,
			st.Result, st.Character, st.Kills, st.Deaths, st.Assists, matchID, st.PlayerID,
		)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx,
		`UPDATE matches SET status = 'completed', ended_at = now() WHERE match_id = $1`,
		matchID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

type MatchSummary struct {
	MatchID string `json:"match_id"`
	Status  string `json:"status"`
	Map     string `json:"map"`
}

func (s *Store) ListMatchesByStatus(ctx context.Context, status string) ([]MatchSummary, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT match_id, status, map FROM matches WHERE status = $1 ORDER BY created_at DESC`,
		status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MatchSummary
	for rows.Next() {
		var m MatchSummary
		if err := rows.Scan(&m.MatchID, &m.Status, &m.Map); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, rows.Err()
}

type ParticipantSummary struct {
	ParticipantID string `json:"participant_id"`
	PlayerID      string `json:"player_id"`
	Team          string `json:"team"`
}

func (s *Store) GetParticipants(ctx context.Context, matchID string) ([]ParticipantSummary, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT participant_id, player_id, team FROM match_participants WHERE match_id = $1`,
		matchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ParticipantSummary
	for rows.Next() {
		var p ParticipantSummary
		if err := rows.Scan(&p.ParticipantID, &p.PlayerID, &p.Team); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, rows.Err()
}
