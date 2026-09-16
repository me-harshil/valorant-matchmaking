package shared

type Team string

const (
	TeamA Team = "A"
	TeamB Team = "B"
)

type Result string

const (
	ResultWin  Result = "win"
	ResultLoss Result = "loss"
	ResultDraw Result = "draw"
)

type MatchParticipant struct {
	ParticipantID    string   `json:participant_id`
	MatchID          string   `json:match_id`
	PlayerID         string   `json:player_id`
	Team             Team     `json:team`
	Result           Result   `json:result`
	Character        string   `json:character`
	Kills            int      `json:kills`
	Deaths           int      `json:deaths`
	Assists          int      `json:assists`
	PerformanceScore *float64 `json:performance_score,omitempty`
}
