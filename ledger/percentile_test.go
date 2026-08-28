package ledger

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// tolerance for comparing the Go result against the generated reference value.
const pctTolerance = 1e-9

// minCorpusCases mirrors MIN_CASES in gen/gen_pct.py.
const minCorpusCases = 60

type pctCase struct {
	Name     string  `json:"name"`
	Xs       []int   `json:"xs"`
	P        float64 `json:"p"`
	Expected float64 `json:"expected"`
}

type pctCorpus struct {
	Generator string    `json:"generator"`
	Method    string    `json:"method"`
	Seed      int       `json:"seed"`
	Cases     []pctCase `json:"cases"`
}

func loadCorpus(t *testing.T) pctCorpus {
	t.Helper()

	path := filepath.Join("..", "data", "pct_cases.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v (regenerate it with: python3 gen/gen_pct.py)", path, err)
	}

	var corpus pctCorpus
	if err := json.Unmarshal(raw, &corpus); err != nil {
		t.Fatalf("decoding %s: %v (regenerate it with: python3 gen/gen_pct.py)", path, err)
	}
	if len(corpus.Cases) < minCorpusCases {
		t.Fatalf("%s has %d cases, want at least %d (regenerate it with: python3 gen/gen_pct.py)",
			path, len(corpus.Cases), minCorpusCases)
	}
	return corpus
}

func TestPercentileCorpus(t *testing.T) {
	corpus := loadCorpus(t)

	for _, c := range corpus.Cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			got := Percentile(c.Xs, c.P)
			if diff := math.Abs(got - c.Expected); diff > pctTolerance {
				t.Errorf("Percentile(%v, %v) = %v, want %v (diff %v)", c.Xs, c.P, got, c.Expected, diff)
			}
		})
	}
}

func TestPercentileEmpty(t *testing.T) {
	if got := Percentile(nil, 50); got != 0 {
		t.Errorf("Percentile(nil, 50) = %v, want 0", got)
	}
	if got := Percentile([]int{}, 50); got != 0 {
		t.Errorf("Percentile([]int{}, 50) = %v, want 0", got)
	}
}

func TestPercentileSingle(t *testing.T) {
	for _, p := range []float64{0, 50, 100} {
		if got := Percentile([]int{7}, p); got != 7 {
			t.Errorf("Percentile([]int{7}, %v) = %v, want 7", p, got)
		}
	}
}

func TestPercentileDoesNotMutateInput(t *testing.T) {
	xs := []int{9, 1, 8, 2, 7, 3}
	before := append([]int(nil), xs...)

	Percentile(xs, 40)

	if !reflect.DeepEqual(xs, before) {
		t.Errorf("Percentile mutated its input: got %v, want %v", xs, before)
	}
}

func TestPercentileClampsP(t *testing.T) {
	xs := []int{4, 8, 15, 16, 23, 42}

	if got, want := Percentile(xs, -10), Percentile(xs, 0); got != want {
		t.Errorf("Percentile(xs, -10) = %v, want %v (clamped to p=0)", got, want)
	}
	if got, want := Percentile(xs, 250), Percentile(xs, 100); got != want {
		t.Errorf("Percentile(xs, 250) = %v, want %v (clamped to p=100)", got, want)
	}
	if got, want := Percentile(xs, math.NaN()), Percentile(xs, 0); got != want {
		t.Errorf("Percentile(xs, NaN) = %v, want %v (treated as p=0)", got, want)
	}
}

func TestPercentileUnsortedMatchesSorted(t *testing.T) {
	unsorted := []int{42, 4, 23, 8, 16, 15}
	sorted := []int{4, 8, 15, 16, 23, 42}

	for _, p := range []float64{0, 12.5, 25, 50, 75, 99, 100} {
		if got, want := Percentile(unsorted, p), Percentile(sorted, p); got != want {
			t.Errorf("Percentile(unsorted, %v) = %v, want %v", p, got, want)
		}
	}
}

func TestPercentileKnownValues(t *testing.T) {
	xs := []int{1, 2, 3, 4}
	for _, tc := range []struct {
		p    float64
		want float64
	}{
		{0, 1},
		{25, 1.75},
		{50, 2.5},
		{75, 3.25},
		{100, 4},
	} {
		if got := Percentile(xs, tc.p); math.Abs(got-tc.want) > pctTolerance {
			t.Errorf("Percentile(%v, %v) = %v, want %v", xs, tc.p, got, tc.want)
		}
	}
}
