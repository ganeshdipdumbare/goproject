package ledger

import (
	"encoding/json"
	"math"
	"os"
	"reflect"
	"testing"
)

// tolerance for comparing the Go result against the reference values in the
// generated corpus. Both implementations use the identical algorithm, so the
// only expected difference is floating-point rounding.
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
		t.Fatalf("corpus has %d cases, want at least 60", len(cases))
	}
	return cases
}

func TestPercentileCorpus(t *testing.T) {
	for i, c := range loadPctCases(t) {
		// Keep an untouched copy to assert non-mutation below.
		orig := append([]int(nil), c.Xs...)

		got := Percentile(c.Xs, c.P)
		if math.Abs(got-c.Expected) > pctTolerance {
			t.Errorf("case %d: Percentile(%v, %g) = %v, want %v", i, c.Xs, c.P, got, c.Expected)
		}

		if len(c.Xs) > 0 && !reflect.DeepEqual(c.Xs, orig) {
			t.Errorf("case %d: Percentile mutated its input: got %v, want %v", i, c.Xs, orig)
		}
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

func TestPercentileSingle(t *testing.T) {
	for _, p := range []float64{0, 25, 50, 100} {
		if got := Percentile([]int{7}, p); got != 7 {
			t.Errorf("Percentile([7], %g) = %v, want 7", p, got)
		}
	}
}

func TestPercentileDoesNotMutate(t *testing.T) {
	xs := []int{5, 3, 9, 1, 7}
	orig := append([]int(nil), xs...)
	_ = Percentile(xs, 50)
	if !reflect.DeepEqual(xs, orig) {
		t.Errorf("Percentile mutated its input: got %v, want %v", xs, orig)
	}
}
