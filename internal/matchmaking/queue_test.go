package matchmaking

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestQueueConcurrentAdds(t *testing.T) {
	q := NewQueue()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			q.AddPlayer(QueuedPlayer{
				PlayerID:           "player-" + strconv.Itoa(n),
				ConservativeRating: float64(n),
				QueuedAt:           time.Now(),
			})
		}(i)
	}
	wg.Wait()

	snap := q.Snapshot()
	if len(snap) != 100 {
		t.Errorf("got %d players in queue, want 100", len(snap))
	}
}

func TestQueueRemove(t *testing.T) {
	q := NewQueue()
	q.AddPlayer(QueuedPlayer{PlayerID: "p1", ConservativeRating: 25, QueuedAt: time.Now()})

	removed := q.RemovePlayer("p1")
	if !removed {
		t.Error("expected RemovePlayer to return true for existing player")
	}

	removedAgain := q.RemovePlayer("p1")
	if removedAgain {
		t.Error("expected RemovePlayer to return false for already-removed player")
	}
}
