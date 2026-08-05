package ledger

import (
	"math"
	"sort"
)

// Percentile returns the p-th percentile of xs using linear interpolation
// between the two closest ranks (the numpy "linear"/type-7 method). p is
// expressed on a 0..100 scale and is clamped into that range.
//
// The function is safe on empty input, returning 0 when len(xs) == 0. It never
// mutates the caller's slice: the values are copied into a fresh slice before
// sorting.
func Percentile(xs []int, p float64) float64 {
	n := len(xs)
	if n == 0 {
		return 0
	}

	// Copy so the caller's slice order is preserved.
	sorted := make([]int, n)
	copy(sorted, xs)
	sort.Ints(sorted)

	if n == 1 {
		return float64(sorted[0])
	}

	// Clamp p to [0, 100] to keep the function total.
	if p < 0 {
		p = 0
	} else if p > 100 {
		p = 100
	}

	rank := (p / 100.0) * float64(n-1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	frac := rank - float64(lo)

	return float64(sorted[lo]) + frac*float64(sorted[hi]-sorted[lo])
}
