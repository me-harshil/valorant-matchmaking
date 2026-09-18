package matchmaking

import (
	"testing"
	"time"
)

func makePlayer(id string, rating float64, secondsAgo float64) QueuedPlayer {
	return QueuedPlayer{
		PlayerID:           id,
		ConservativeRating: rating,
		QueuedAt:           time.Now().Add(-time.Duration(secondsAgo) * time.Second),
	}
}

func TestFindMatch_TightClusterMatchesImmediately(t *testing.T) {
	var players []QueuedPlayer
	for i := 0; i < 10; i++ {
		// ratings 40..49, all just joined -> spread=9, well within base window of 50
		players = append(players, makePlayer("p"+string(rune('a'+i)), float64(40+i), 0))
	}

	now := float64(time.Now().Unix())
	match := FindMatch(players, now)

	if match == nil {
		t.Fatal("expected a match, got nil")
	}
	if len(match) != 10 {
		t.Errorf("got %d players, want 10", len(match))
	}
}

func TestFindMatch_TooFewPlayersReturnsNil(t *testing.T) {
	players := []QueuedPlayer{makePlayer("p1", 25, 0)}
	match := FindMatch(players, float64(time.Now().Unix()))
	if match != nil {
		t.Error("expected nil with fewer than 10 players")
	}
}

func TestFindMatch_WideSpreadRejectedWhenFreshlyQueued(t *testing.T) {
	var players []QueuedPlayer
	for i := 0; i < 10; i++ {
		// ratings spread across 0..900 -> spread way beyond base window, all just joined
		players = append(players, makePlayer("p"+string(rune('a'+i)), float64(i*100), 0))
	}
	now := float64(time.Now().Unix())
	match := FindMatch(players, now)
	if match != nil {
		t.Error("expected nil — spread too wide for freshly-queued players")
	}
}
