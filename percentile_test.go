package ledger

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"
)

type pctCase struct {
	Xs       []int   `json:"xs"`
	P        float64 `json:"p"`
	Expected float64 `json:"expected"`
}

const pctTolerance = 1e-9

func TestPercentileCorpus(t *testing.T) {
	path := filepath.Join("data", "pct_cases.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	var cases []pctCase
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatalf("unmarshaling %s: %v", path, err)
	}
	if len(cases) < 60 {
		t.Fatalf("expected at least 60 cases, got %d", len(cases))
	}

	for i, c := range cases {
		got := Percentile(c.Xs, c.P)
		if math.Abs(got-c.Expected) > pctTolerance {
			t.Errorf("case %d: Percentile(%v, %v) = %v, want %v", i, c.Xs, c.P, got, c.Expected)
		}
	}
}

func TestPercentileEmpty(t *testing.T) {
	if got := Percentile(nil, 50); got != 0.0 {
		t.Errorf("Percentile(nil, 50) = %v, want 0", got)
	}
	if got := Percentile([]int{}, 99); got != 0.0 {
		t.Errorf("Percentile([]int{}, 99) = %v, want 0", got)
	}
}

func TestPercentileDoesNotMutate(t *testing.T) {
	xs := []int{5, 1, 4, 2, 3}
	orig := make([]int, len(xs))
	copy(orig, xs)

	_ = Percentile(xs, 50)

	for i := range xs {
		if xs[i] != orig[i] {
			t.Fatalf("input slice mutated: got %v, want %v", xs, orig)
		}
	}
}

func TestPercentileBoundaries(t *testing.T) {
	xs := []int{10, 20, 30, 40}
	if got := Percentile(xs, 0); got != 10 {
		t.Errorf("Percentile(p=0) = %v, want 10", got)
	}
	if got := Percentile(xs, 100); got != 40 {
		t.Errorf("Percentile(p=100) = %v, want 40", got)
	}
}
