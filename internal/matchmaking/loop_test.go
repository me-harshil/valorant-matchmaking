package matchmaking

import (
	"testing"
	"time"
)

func TestMatcherLoop(t *testing.T) {
	q := NewQueue()
	for i := 0; i < 10; i++ {
		q.AddPlayer(QueuedPlayer{
			PlayerID:           "p" + string(rune('a'+i)),
			ConservativeRating: float64(40 + i),
			QueuedAt:           time.Now(),
		})
	}

	matched := make(chan bool, 1)

	go RunMatcherLoop(q, 50*time.Millisecond, func(teamA, teamB []QueuedPlayer) {
		if len(teamA) != 5 || len(teamB) != 5 {
			t.Errorf("expected two teams of 5 players each, got %d and %d", len(teamA), len(teamB))
		}
		matched <- true
	})

	select {
	case <-matched:
		// success - callback fired
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for a match")
	}
}
