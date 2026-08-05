package ledger

import (
	"math"
	"sort"
)

// Percentile returns the p-th percentile of xs using linear interpolation
// between the two closest ranks (the same "linear" method NumPy uses by
// default). p is expressed on a 0..100 scale and is clamped into that range.
//
// The function is safe on empty input: it returns 0 when xs is empty. It never
// mutates the caller's slice, operating instead on an internal sorted copy.
func Percentile(xs []int, p float64) float64 {
	n := len(xs)
	if n == 0 {
		return 0
	}

	sorted := make([]int, n)
	copy(sorted, xs)
	sort.Ints(sorted)

	if n == 1 {
		return float64(sorted[0])
	}

	if p < 0 {
		p = 0
	}
	if p > 100 {
		p = 100
	}

	rank := (p / 100) * float64(n-1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	frac := rank - float64(lo)

	return float64(sorted[lo]) + frac*(float64(sorted[hi])-float64(sorted[lo]))
}
