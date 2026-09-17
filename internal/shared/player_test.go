package shared

import "testing"

func TestConservativeRating(t *testing.T) {
	p := Player{Mu: 25.0, Sigma: 8.333}

	got := p.ConservativeRating()
	want := 25.0 - 3*8.333 // 0.001

	const epsilon = 0.0001
	diff := got - want
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Errorf("ConservativeRating() = %v, want %v", got, want)
	}
}
