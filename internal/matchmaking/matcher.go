package matchmaking

import (
	"math/bits"
	"sort"
)

const (
	baseWindow       = 50.0
	windowGrowthRate = 50.0 // window widens by this much per growth interval
	growthInterval   = 15.0 // seconds waited before window widens once
	playersPerMatch  = 10
	teamSize         = 5
)

func currentWindow(waitedSeconds float64) float64 {
	growthSteps := waitedSeconds / growthInterval
	return baseWindow + growthSteps*windowGrowthRate
}

func FindMatch(players []QueuedPlayer, nowUnixSeconds float64) []QueuedPlayer {
	if len(players) < playersPerMatch {
		return nil
	}

	sorted := make([]QueuedPlayer, len(players))
	copy(sorted, players)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ConservativeRating < sorted[j].ConservativeRating
	})

	for i := 0; i+playersPerMatch <= len(sorted); i++ {
		group := sorted[i : i+playersPerMatch]
		if groupIsCompatible(group, nowUnixSeconds) {
			return group
		}
	}
	return nil
}

func groupIsCompatible(group []QueuedPlayer, nowUnixSeconds float64) bool {
	spread := group[len(group)-1].ConservativeRating - group[0].ConservativeRating

	for _, p := range group {
		waited := nowUnixSeconds - float64(p.QueuedAt.Unix())
		if currentWindow(waited) < spread {
			return false
		}
	}
	return true
}

// SplitTeams finds the 5v5 split of exactly 10 players that minimizes the difference between team average ratings
func SplitTeams(players []QueuedPlayer) (teamA, teamB []QueuedPlayer) {
	bestDiff := -1.0
	var bestA, bestB []QueuedPlayer

	// iterate all C(10,5) = 252 combinations
	n := len(players)

	for mask := 0; mask < (1 << n); mask++ {
		if bits.OnesCount(uint(mask)) != teamSize {
			continue
		}

		var a, b []QueuedPlayer
		var sumA, sumB float64
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				a = append(a, players[i])
				sumA += players[i].ConservativeRating
			} else {
				b = append(b, players[i])
				sumB += players[i].ConservativeRating
			}
		}

		avgA := sumA / float64(len(a))
		avgB := sumB / float64(len(b))
		diff := avgA - avgB
		if diff < 0 {
			diff = -diff
		}

		if bestDiff < 0 || diff < bestDiff {
			bestDiff = diff
			bestA, bestB = a, b
		}
	}

	return bestA, bestB
}
