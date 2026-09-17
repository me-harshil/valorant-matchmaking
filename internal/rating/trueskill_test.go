package rating

import "testing"

func almostEqual(a, b, epsilon float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= epsilon
}

func TestUpdateRatingsTwoOnTwo(t *testing.T) {
	teamA := []PlayerInput{{"p1", 25.0, 8.333}, {"p2", 25.0, 8.333}}
	teamB := []PlayerInput{{"p3", 25.0, 8.333}, {"p4", 25.0, 8.333}}

	newA, _ := UpdateRatings(teamA, teamB, true, false)

	if !almostEqual(newA[0].Mu, 28.108, 0.085) {
		t.Errorf("Mu = %v, want ~28.108", newA[0].Mu)
	}
	if !almostEqual(newA[0].Sigma, 7.774, 0.085) {
		t.Errorf("Sigma = %v, want ~7.774", newA[0].Sigma)
	}
}
