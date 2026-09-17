package rating

import "testing"

func TestApplyPerformanceWeight_MultiplierOne(t *testing.T) {
	got := ApplyPerformanceWeight(25.0, 28.108, 1.0)
	if !almostEqual(got, 28.108, 0.0001) {
		t.Errorf("got %v, want 28.108 (multiplier=1 should be a no-op)", got)
	}
}

func TestNormalizeWithinTeam_Clamping(t *testing.T) {
	team := []PerformanceStats{
		{"p1", 20, 1, 5}, // big performer
		{"p2", 1, 10, 0}, // very poor performer
	}
	result := NormalizeWithinTeam(team)

	if result["p1"] > 1.5 {
		t.Errorf("p1 multiplier %v exceeds clamp of 1.5", result["p1"])
	}
	if result["p2"] < 0.5 {
		t.Errorf("p2 multiplier %v below clamp of 0.5", result["p2"])
	}
}
