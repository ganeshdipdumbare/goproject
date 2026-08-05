package ledger

import (
	"math"
	"sort"
)

// Percentile returns the p-th percentile of xs using linear interpolation
// between the two closest ranks, matching numpy's default "linear" method.
//
// The values are treated as a sample sorted ascending. For a slice of length
// n the fractional rank is r = (n-1) * (p/100) on the 0-indexed sorted data;
// the result interpolates linearly between the values at floor(r) and ceil(r).
//
// p is clamped to the range [0, 100]. An empty slice returns 0. The caller's
// slice is never mutated: Percentile sorts a copy.
func Percentile(xs []int, p float64) float64 {
	if len(xs) == 0 {
		return 0.0
	}
	if p < 0 {
		p = 0
	}
	if p > 100 {
		p = 100
	}

	sorted := make([]int, len(xs))
	copy(sorted, xs)
	sort.Ints(sorted)

	if len(sorted) == 1 {
		return float64(sorted[0])
	}

	r := float64(len(sorted)-1) * (p / 100.0)
	lo := int(math.Floor(r))
	hi := int(math.Ceil(r))
	frac := r - float64(lo)

	return float64(sorted[lo]) + frac*float64(sorted[hi]-sorted[lo])
}
