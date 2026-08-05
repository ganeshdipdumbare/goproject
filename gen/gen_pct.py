#!/usr/bin/env python3
"""Deterministic generator for the Percentile test corpus.

Emits data/pct_cases.json containing at least 60 seeded random cases. Each case
records an input slice, a percentile p (0..100), and the expected value computed
by the reference implementation below. The reference matches the Go
implementation in percentile.go exactly (NumPy's default "linear" method), so
the committed JSON is byte-stable across runs.
"""

import json
import math
import os
import random


def percentile(xs, p):
    """Reference percentile using linear interpolation between closest ranks.

    Mirrors ledger.Percentile in percentile.go:
      * empty input returns 0
      * never mutates the caller's slice (operates on a sorted copy)
      * p is clamped into [0, 100]
    """
    n = len(xs)
    if n == 0:
        return 0.0

    s = sorted(xs)
    if n == 1:
        return float(s[0])

    p = max(0.0, min(100.0, p))
    rank = (p / 100.0) * (n - 1)
    lo = math.floor(rank)
    hi = math.ceil(rank)
    frac = rank - lo
    return float(s[lo]) + frac * (float(s[hi]) - float(s[lo]))


def main():
    rng = random.Random(20240115)
    cases = []

    # A few explicit edge cases up front for readability.
    cases.append({"input": [], "p": 50.0})
    cases.append({"input": [42], "p": 50.0})
    cases.append({"input": [1, 2, 3, 4, 5], "p": 0.0})
    cases.append({"input": [1, 2, 3, 4, 5], "p": 100.0})

    # Generated seeded random cases spanning varied lengths and p values.
    for _ in range(70):
        length = rng.randint(1, 25)
        xs = [rng.randint(-100, 100) for _ in range(length)]
        p = round(rng.uniform(0.0, 100.0), 4)
        cases.append({"input": xs, "p": p})

    for c in cases:
        c["expected"] = percentile(c["input"], c["p"])

    out_dir = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "data")
    os.makedirs(out_dir, exist_ok=True)
    out_path = os.path.join(out_dir, "pct_cases.json")

    with open(out_path, "w") as f:
        json.dump(cases, f, indent=2, sort_keys=True)
        f.write("\n")


if __name__ == "__main__":
    main()
