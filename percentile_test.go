package ledger

import (
	"encoding/json"
	"math"
	"os"
	"testing"
)

// epsilon tolerates last-ULP differences between Python's and Go's
// floating-point formatting when comparing generated expectations.
const epsilon = 1e-9

type pctCase struct {
	Xs       []int   `json:"xs"`
	P        float64 `json:"p"`
	Expected float64 `json:"expected"`
}

func loadCases(t *testing.T) []pctCase {
	t.Helper()

	data, err := os.ReadFile("data/pct_cases.json")
	if err != nil {
		t.Fatalf("reading data/pct_cases.json: %v", err)
	}

	var cases []pctCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshaling data/pct_cases.json: %v", err)
	}
	if len(cases) < 60 {
		t.Fatalf("expected at least 60 generated cases, got %d", len(cases))
	}
	return cases
}

func TestPercentileGeneratedCases(t *testing.T) {
	cases := loadCases(t)

	for i, c := range cases {
		got := Percentile(c.Xs, c.P)
		if math.Abs(got-c.Expected) > epsilon {
			t.Errorf("case %d: Percentile(%v, %v) = %v, want %v",
				i, c.Xs, c.P, got, c.Expected)
		}
	}
}

func TestPercentileEmpty(t *testing.T) {
	if got := Percentile(nil, 0.5); got != 0 {
		t.Errorf("Percentile(nil, 0.5) = %v, want 0", got)
	}
	if got := Percentile([]int{}, 0.9); got != 0 {
		t.Errorf("Percentile([]int{}, 0.9) = %v, want 0", got)
	}
}

func TestPercentileDoesNotMutateInput(t *testing.T) {
	xs := []int{5, 1, 4, 2, 3}
	before := make([]int, len(xs))
	copy(before, xs)

	_ = Percentile(xs, 0.5)

	for i := range xs {
		if xs[i] != before[i] {
			t.Fatalf("input mutated: got %v, want %v", xs, before)
		}
	}
}
