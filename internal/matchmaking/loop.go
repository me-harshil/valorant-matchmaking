package matchmaking

import (
	"log"
	"time"
)

// RunMatcherLoop ticks every interval, checks the queue for a match, and if found, removes those 10 players and calls onMatch with the two teams
func RunMatcherLoop(q *Queue, interval time.Duration, onMatch func(teamA, teamB []QueuedPlayer) error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		snapshot := q.Snapshot()
		now := float64(time.Now().Unix())

		match := FindMatch(snapshot, now)
		if match == nil {
			continue
		}

		for _, p := range match {
			q.RemovePlayer(p.PlayerID)
		}

		teamA, teamB := SplitTeams(match)
		if err := onMatch(teamA, teamB); err != nil {
			log.Printf("match creation failed, requeueing %d players: %v", len(match), err)
			for _, p := range match {
				q.AddPlayer(p)
			}
		}
	}
}
