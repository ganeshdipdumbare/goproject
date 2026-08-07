package ledger

import (
	"math"
	"sort"
)

// Percentile returns the p-th percentile of xs using linear interpolation
// between the two closest ranks (equivalent to numpy's default "linear"
// method).
//
// The parameter p is a percentile in the range [0, 100]; values outside that
// range are clamped. Empty input returns 0. The caller's slice is never
// mutated: Percentile sorts a private copy of xs.
func Percentile(xs []int, p float64) float64 {
	if len(xs) == 0 {
		return 0
	}

	// Copy so the caller's slice ordering is preserved, then sort ascending.
	sorted := append([]int(nil), xs...)
	sort.Ints(sorted)

	if p < 0 {
		p = 0
	} else if p > 100 {
		p = 100
	}

	n := len(sorted)
	if n == 1 {
		return float64(sorted[0])
	}

	// Fractional rank in [0, n-1].
	r := (p / 100) * float64(n-1)
	lo := int(math.Floor(r))
	hi := int(math.Ceil(r))
	frac := r - float64(lo)

	return float64(sorted[lo]) + frac*(float64(sorted[hi])-float64(sorted[lo]))
}
