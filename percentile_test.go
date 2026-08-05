package ledger

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"
)

const pctTolerance = 1e-9

type percentileCase struct {
	Input    []int   `json:"input"`
	P        float64 `json:"p"`
	Expected float64 `json:"expected"`
}

func TestPercentileCorpus(t *testing.T) {
	data, err := os.ReadFile("data/pct_cases.json")
	if err != nil {
		t.Fatalf("reading corpus: %v", err)
	}

	var cases []percentileCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("decoding corpus: %v", err)
	}

	if len(cases) < 60 {
		t.Fatalf("expected at least 60 cases, got %d", len(cases))
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			got := Percentile(c.Input, c.P)
			if math.Abs(got-c.Expected) > pctTolerance {
				t.Errorf("Percentile(%v, %v) = %v, want %v", c.Input, c.P, got, c.Expected)
			}
		})
	}
}

func TestPercentileEmpty(t *testing.T) {
	if got := Percentile(nil, 50); got != 0 {
		t.Errorf("Percentile(nil, 50) = %v, want 0", got)
	}
	if got := Percentile([]int{}, 25); got != 0 {
		t.Errorf("Percentile([], 25) = %v, want 0", got)
	}
}

func TestPercentileSingle(t *testing.T) {
	for _, p := range []float64{0, 25, 50, 100} {
		if got := Percentile([]int{7}, p); got != 7 {
			t.Errorf("Percentile([7], %v) = %v, want 7", p, got)
		}
	}
}

func TestPercentileBoundaries(t *testing.T) {
	xs := []int{5, 1, 3, 2, 4}
	if got := Percentile(xs, 0); got != 1 {
		t.Errorf("Percentile(%v, 0) = %v, want 1", xs, got)
	}
	if got := Percentile(xs, 100); got != 5 {
		t.Errorf("Percentile(%v, 100) = %v, want 5", xs, got)
	}
	if got := Percentile(xs, 50); got != 3 {
		t.Errorf("Percentile(%v, 50) = %v, want 3", xs, got)
	}
}

func TestPercentileDoesNotMutate(t *testing.T) {
	xs := []int{5, 1, 3, 2, 4}
	before := make([]int, len(xs))
	copy(before, xs)

	_ = Percentile(xs, 42)

	if !reflect.DeepEqual(xs, before) {
		t.Errorf("Percentile mutated input: got %v, want %v", xs, before)
	}
}
