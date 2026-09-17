package rating

import (
	"github.com/ChrisHines/GoSkills/skills"
	"github.com/ChrisHines/GoSkills/skills/trueskill"
)

type PlayerInput struct {
	PlayerID      string
	ParticipantID string
	Mu            float64
	Sigma         float64
}

type PlayerOutput struct {
	PlayerID      string
	ParticipantID string
	Mu            float64
	Sigma         float64
}

func UpdateRatings(teamA, teamB []PlayerInput, teamAWon bool, isDraw bool) (newA, newB []PlayerOutput) {
	gi := skills.DefaultGameInfo
	calc := &trueskill.TwoTeamCalc{}

	skTeamA := skills.NewTeam()
	for _, p := range teamA {
		skTeamA.AddPlayer(p.PlayerID, skills.NewRating(p.Mu, p.Sigma))
	}

	skTeamB := skills.NewTeam()
	for _, p := range teamB {
		skTeamB.AddPlayer(p.PlayerID, skills.NewRating(p.Mu, p.Sigma))
	}

	ranks := []int{1, 2}
	if isDraw {
		ranks = []int{1, 1}
	} else if !teamAWon {
		ranks = []int{2, 1}
	}

	result := calc.CalcNewRatings(gi, []skills.Team{skTeamA, skTeamB}, ranks...)

	for _, p := range teamA {
		r := result[p.PlayerID]
		newA = append(newA, PlayerOutput{
			PlayerID:      p.PlayerID,
			ParticipantID: p.ParticipantID,
			Mu:            r.Mean(),
			Sigma:         r.Stddev(),
		})
	}

	for _, p := range teamB {
		r := result[p.PlayerID]
		newB = append(newB, PlayerOutput{
			PlayerID:      p.PlayerID,
			ParticipantID: p.ParticipantID,
			Mu:            r.Mean(),
			Sigma:         r.Stddev(),
		})
	}
	return
}
