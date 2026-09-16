package shared

import "time"

type MatchStatus string

const (
	MatchStatusPending    MatchStatus = "pending"
	MatchStatusInProgress MatchStatus = "in_progress"
	MatchStatusCompleted  MatchStatus = "completed"
	MatchStatusCancelled  MatchStatus = "cancelled"
)

type Match struct {
	MatchID   string      `json:match_id`
	Status    MatchStatus `json:status`
	Map       string      `json:map`
	StartedAt *time.Time  `json:started_at,omitempty`
	EndedAt   *time.Time  `json:ended_at,omitempty`
	CreatedAt time.Time   `json:created_at`
}
