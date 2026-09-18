package matchmaking

import "sort"

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
