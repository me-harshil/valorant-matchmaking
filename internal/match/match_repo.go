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
