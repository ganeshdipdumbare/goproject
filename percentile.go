package ledger

import (
	"math"
	"sort"
)

// Percentile returns the p-th percentile of xs using linear interpolation
// between closest ranks (the same definition as numpy.percentile with
// method="linear").
//
// p is interpreted on the domain [0, 1] and is clamped to that range, so
// p=0 yields the minimum, p=1 the maximum, and p=0.5 the median.
//
// The function never mutates the caller's slice: it sorts an internal copy.
// For empty input it returns 0 as a documented sentinel rather than panicking.
func Percentile(xs []int, p float64) float64 {
	n := len(xs)
	if n == 0 {
		return 0
	}

	// Copy so the caller's slice is never mutated, then sort ascending.
	sorted := make([]int, n)
	copy(sorted, xs)
	sort.Ints(sorted)

	if n == 1 {
		return float64(sorted[0])
	}

	// Clamp p to [0, 1].
	if p < 0 {
		p = 0
	} else if p > 1 {
		p = 1
	}

	// Fractional rank across the sorted copy.
	r := p * float64(n-1)
	lo := int(math.Floor(r))
	hi := int(math.Ceil(r))
	frac := r - float64(lo)

	if lo == hi {
		return float64(sorted[lo])
	}
	return float64(sorted[lo]) + frac*float64(sorted[hi]-sorted[lo])
}
