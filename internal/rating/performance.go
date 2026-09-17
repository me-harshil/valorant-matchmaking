package rating

type PerformanceStats struct {
	PlayerID string
	Kills    int
	Deaths   int
	Assists  int
}

const (
	weightKill   = 1.0
	weightDeath  = -1.0
	weightAssist = 0.5
)

func rawScore(s PerformanceStats) float64 {
	return float64(s.Kills)*weightKill + float64(s.Deaths)*weightDeath + float64(s.Assists)*weightAssist
}

// NormalizeWithinTeam returns each player's score divided by the team average, so no single match causes an extreme rating swing
func NormalizeWithinTeam(team []PerformanceStats) map[string]float64 {
	total := 0.0
	for _, s := range team {
		total += rawScore(s)
	}
	avg := total / float64(len(team))

	result := make(map[string]float64)

	for _, s := range team {
		mult := 1.0
		if avg != 0 {
			mult = rawScore(s) / avg
		}
		if mult < 0.5 {
			mult = 0.5
		}
		if mult > 1.5 {
			mult = 1.5
		}
		result[s.PlayerID] = mult
	}
	return result
}

func ApplyPerformanceWeight(oldMu, rawNewMu, multiplier float64) float64 {
	delta := rawNewMu - oldMu
	return oldMu + delta*multiplier
}
