package ledger

import (
	"math"
	"sort"
)

// Percentile returns the p-th percentile of xs using linear interpolation
// between closest ranks (equivalent to NumPy's default method="linear", i.e.
// C=1).
//
// The percentile p is clamped to the [0, 100] domain. The computation copies
// xs before sorting, so the caller's slice is never mutated and its ordering
// is preserved.
//
// Edge cases:
//   - an empty slice returns 0.
//   - a single-element slice returns that element for any p.
//
// The rank is computed as rank = (p/100) * (n-1); with lo = floor(rank),
// hi = ceil(rank) and frac = rank - lo the result is
// sorted[lo] + frac*(sorted[hi]-sorted[lo]).
func Percentile(xs []int, p float64) float64 {
	n := len(xs)
	if n == 0 {
		return 0
	}

	// Copy so we never mutate the caller's slice.
	sorted := make([]int, n)
	copy(sorted, xs)
	sort.Ints(sorted)

	if n == 1 {
		return float64(sorted[0])
	}

	// Clamp p to the [0, 100] domain so Go and the reference generator agree.
	if p < 0 {
		p = 0
	} else if p > 100 {
		p = 100
	}

	rank := (p / 100) * float64(n-1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	frac := rank - float64(lo)

	return float64(sorted[lo]) + frac*float64(sorted[hi]-sorted[lo])
}
