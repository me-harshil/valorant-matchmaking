package matchmaking

import "time"

type QueuedPlayer struct {
	PlayerID           string
	ConservativeRating float64
	QueuedAt           time.Time
}

// addRequest is sent to the queue-owning goroutine to add a player.
type addRequest struct {
	player QueuedPlayer
}

// removeRequest is sent to remove a player (e.g. they cancel queueing)
type removeRequest struct {
	playerID string
	done     chan bool // true if actually removed, false if not found
}

type snapshotRequest struct {
	result chan []QueuedPlayer
}

type Queue struct {
	add      chan addRequest
	remove   chan removeRequest
	snapshot chan snapshotRequest
}

func NewQueue() *Queue {
	q := &Queue{
		add:      make(chan addRequest),
		remove:   make(chan removeRequest),
		snapshot: make(chan snapshotRequest),
	}
	go q.run()
	return q
}

func (q *Queue) run() {
	queue := make(map[string]QueuedPlayer)

	for {
		select {
		case req := <-q.add:
			queue[req.player.PlayerID] = req.player

		case req := <-q.remove:
			_, existed := queue[req.playerID]
			delete(queue, req.playerID)
			req.done <- existed

		case req := <-q.snapshot:
			cp := make([]QueuedPlayer, 0, len(queue))
			for _, p := range queue {
				cp = append(cp, p)
			}
			req.result <- cp
		}
	}
}

func (q *Queue) AddPlayer(p QueuedPlayer) {
	q.add <- addRequest{player: p}
}

func (q *Queue) RemovePlayer(playerID string) bool {
	done := make(chan bool)
	q.remove <- removeRequest{playerID: playerID, done: done}
	return <-done
}

func (q *Queue) Snapshot() []QueuedPlayer {
	result := make(chan []QueuedPlayer)
	q.snapshot <- snapshotRequest{result: result}
	return <-result
}
