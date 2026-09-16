package shared

import "time"

type Player struct {
	PlayerID  string    `json:player_id`
	Username  string    `json:username`
	Mu        float64   `json:mu`
	Sigma     float64   `json:sigma`
	CreatedAt time.Time `json:created_at`
	UpdatedAt time.Time `json:updated_at`
}

func (p *Player) ConservativeRating() float64 {
	return p.Mu - 3*p.Sigma
}
