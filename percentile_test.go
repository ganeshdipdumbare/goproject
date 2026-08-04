package ledger

import (
	"encoding/json"
	"math"
	"os"
	"testing"
)

const pctTolerance = 1e-9

type pctCase struct {
	Xs       []int   `json:"xs"`
	P        float64 `json:"p"`
	Expected float64 `json:"expected"`
}

func loadPctCases(t *testing.T) []pctCase {
	t.Helper()
	data, err := os.ReadFile("data/pct_cases.json")
	if err != nil {
		t.Fatalf("reading corpus: %v", err)
	}
	var cases []pctCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("parsing corpus: %v", err)
	}
	if len(cases) < 60 {
		t.Fatalf("expected at least 60 cases, got %d", len(cases))
	}
	return cases
}

func TestPercentileCorpus(t *testing.T) {
	for i, c := range loadPctCases(t) {
		got := Percentile(c.Xs, c.P)
		if math.Abs(got-c.Expected) > pctTolerance {
			t.Errorf("case %d: Percentile(%v, %v) = %v, want %v",
				i, c.Xs, c.P, got, c.Expected)
		}
	}
}

func TestPercentileEmpty(t *testing.T) {
	for _, p := range []float64{0, 25, 50, 100} {
		if got := Percentile(nil, p); got != 0 {
			t.Errorf("Percentile(nil, %v) = %v, want 0", p, got)
		}
		if got := Percentile([]int{}, p); got != 0 {
			t.Errorf("Percentile([], %v) = %v, want 0", p, got)
		}
	}
}

func TestPercentileSingleElement(t *testing.T) {
	for _, p := range []float64{0, 25, 50, 100} {
		if got := Percentile([]int{7}, p); got != 7 {
			t.Errorf("Percentile([7], %v) = %v, want 7", p, got)
		}
	}
}

func TestPercentileDoesNotMutate(t *testing.T) {
	xs := []int{4, 1, 3, 2}
	orig := make([]int, len(xs))
	copy(orig, xs)

	_ = Percentile(xs, 50)

	for i := range xs {
		if xs[i] != orig[i] {
			t.Fatalf("input mutated: got %v, want %v", xs, orig)
		}
	}
}
