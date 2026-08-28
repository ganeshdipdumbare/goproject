#!/usr/bin/env python3
"""Generate the percentile reference corpus consumed by the Go tests.

The corpus lives in ``data/pct_cases.json`` and must always be exactly what
this script emits. Re-run it from anywhere after editing this file::

    python3 gen/gen_pct.py

then check that ``git diff data/pct_cases.json`` is clean (the output is
deterministic: a fixed seed and stable JSON formatting).

The reference implementation below mirrors ``Percentile`` in
``ledger/percentile.go`` line for line: linear interpolation between the two closest ranks (the C=1
convention, i.e. ``numpy.percentile(..., method="linear")``). It is written in
pure Python on purpose -- no third-party dependency, and no risk of a numpy
version changing the default interpolation method under us.
"""

import json
import math
import random
from pathlib import Path

SEED = 20240607
MIN_CASES = 60

OUTPUT = Path(__file__).resolve().parent.parent / "data" / "pct_cases.json"

# Slice lengths swept over: singletons, tiny odd/even lengths and a few larger
# ones so both exact-rank and interpolated-rank paths are exercised.
LENGTHS = [1, 2, 3, 4, 5, 7, 8, 10, 50, 101, 500]

# Percentiles swept over: the boundaries, the quartiles, the extremes and a
# few values that never land on an integer rank.
PERCENTILES = [0.0, 1.0, 12.5, 25.0, 33.3, 50.0, 66.6, 75.0, 99.0, 100.0]

# Value ranges cycled through per case: non-negative, spanning zero, strictly
# negative and a narrow band that forces plenty of duplicates.
VALUE_RANGES = [(0, 100), (-50, 50), (-1000, -1), (0, 3), (-10000, 10000)]


def percentile(xs, p):
    """Reference percentile, mirroring Percentile in ledger/percentile.go."""
    if not xs:
        return 0.0

    if math.isnan(p) or p < 0:
        p = 0.0
    elif p > 100:
        p = 100.0

    s = sorted(xs)
    n = len(s)
    if n == 1:
        return float(s[0])

    r = (p / 100) * (n - 1)
    lo = int(math.floor(r))
    hi = int(math.ceil(r))
    if lo == hi:
        return float(s[lo])

    frac = r - lo
    return float(s[lo]) + frac * (float(s[hi]) - float(s[lo]))


def case(name, xs, p):
    return {"name": name, "xs": list(xs), "p": float(p), "expected": percentile(xs, p)}


def build_cases(rng):
    cases = []

    # Deterministic sweep over lengths x percentiles.
    for i, n in enumerate(LENGTHS):
        for j, p in enumerate(PERCENTILES):
            lo, hi = VALUE_RANGES[(i + j) % len(VALUE_RANGES)]
            xs = [rng.randint(lo, hi) for _ in range(n)]
            cases.append(case("sweep_len{}_p{}".format(n, p), xs, p))

    # Hand-picked edge cases with values that are easy to verify by hand.
    cases.append(case("edge_single", [7], 50.0))
    cases.append(case("edge_pair_mid", [1, 2], 50.0))
    cases.append(case("edge_quartile_exact", [1, 2, 3, 4], 25.0))
    cases.append(case("edge_quartile_interp", [1, 2, 3, 4], 50.0))
    cases.append(case("edge_all_equal", [5, 5, 5, 5, 5], 42.0))
    cases.append(case("edge_unsorted", [9, 1, 8, 2, 7, 3], 40.0))
    cases.append(case("edge_negatives", [-5, -4, -3, -2, -1], 10.0))
    cases.append(case("edge_wide_span", [-1000000, 0, 1000000], 75.0))
    cases.append(case("edge_duplicates", [1, 1, 1, 2, 2, 3], 60.0))
    cases.append(case("edge_min", [4, 8, 15, 16, 23, 42], 0.0))
    cases.append(case("edge_max", [4, 8, 15, 16, 23, 42], 100.0))
    cases.append(case("edge_fractional_p", [10, 20, 30, 40, 50], 37.5))

    return cases


def main():
    rng = random.Random(SEED)
    cases = build_cases(rng)

    names = [c["name"] for c in cases]
    assert len(names) == len(set(names)), "case names must be unique"
    assert len(cases) >= MIN_CASES, "expected at least {} cases, got {}".format(
        MIN_CASES, len(cases)
    )

    payload = {
        "generator": "gen/gen_pct.py",
        "method": "linear interpolation between closest ranks",
        "seed": SEED,
        "cases": cases,
    }

    OUTPUT.parent.mkdir(parents=True, exist_ok=True)
    with OUTPUT.open("w", encoding="utf-8") as f:
        json.dump(payload, f, indent=2, sort_keys=True)
        f.write("\n")

    print("wrote {} cases to {}".format(len(cases), OUTPUT))


if __name__ == "__main__":
    main()
