package ledger

import (
	"math"
	"sort"
)

// Percentile returns the p-th percentile of xs using linear interpolation
// between the two closest ranks (the C=1 convention, equivalent to
// numpy.percentile with method="linear"): the requested rank is
// r = (p/100)*(len(xs)-1) and the result interpolates between the order
// statistics at floor(r) and ceil(r).
//
// Empty input is safe and returns 0. p is clamped to [0, 100], and a NaN p is
// treated as 0, so the function never panics or indexes out of range.
//
// The caller's slice is never mutated or reordered: Percentile sorts a copy.
func Percentile(xs []int, p float64) float64 {
	if len(xs) == 0 {
		return 0
	}

	if math.IsNaN(p) || p < 0 {
		p = 0
	} else if p > 100 {
		p = 100
	}

	s := append([]int(nil), xs...)
	sort.Ints(s)

	n := len(s)
	if n == 1 {
		return float64(s[0])
	}

	r := (p / 100) * float64(n-1)
	lo := int(math.Floor(r))
	hi := int(math.Ceil(r))
	if lo == hi {
		return float64(s[lo])
	}

	frac := r - float64(lo)
	return float64(s[lo]) + frac*(float64(s[hi])-float64(s[lo]))
}
