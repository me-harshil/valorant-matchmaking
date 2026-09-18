package matchmaking

import "time"

// RunMatcherLoop ticks every interval, checks the queue for a match, and if found, removes those 10 players and calls onMatch with the two teams
func RunMatcherLoop(q *Queue, interval time.Duration, onMatch func(teamA, teamB []QueuedPlayer)) {
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
		onMatch(teamA, teamB)
	}
}
