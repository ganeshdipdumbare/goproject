package ledger

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
)

const pctEpsilon = 1e-9

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
		t.Fatalf("unmarshaling corpus: %v", err)
	}
	if len(cases) < 60 {
		t.Fatalf("expected at least 60 cases, got %d", len(cases))
	}
	return cases
}

func TestPercentileCorpus(t *testing.T) {
	cases := loadPctCases(t)
	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			got := Percentile(c.Xs, c.P)
			if math.Abs(got-c.Expected) > pctEpsilon {
				t.Errorf("Percentile(%v, %g) = %v, want %v", c.Xs, c.P, got, c.Expected)
			}
		})
	}
}

func TestPercentileEmpty(t *testing.T) {
	if got := Percentile(nil, 50); got != 0 {
		t.Errorf("Percentile(nil, 50) = %v, want 0", got)
	}
	if got := Percentile([]int{}, 99); got != 0 {
		t.Errorf("Percentile([]int{}, 99) = %v, want 0", got)
	}
}

func TestPercentileDoesNotMutate(t *testing.T) {
	xs := []int{5, 1, 4, 2, 3}
	snapshot := make([]int, len(xs))
	copy(snapshot, xs)

	Percentile(xs, 50)

	for i := range xs {
		if xs[i] != snapshot[i] {
			t.Fatalf("input slice mutated: got %v, want %v", xs, snapshot)
		}
	}
}
