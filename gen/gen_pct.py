#!/usr/bin/env python3
"""Deterministic generator for the percentile test corpus.

Emits data/pct_cases.json containing seeded random cases. Each case is an
object {"xs": [...], "p": <float>, "expected": <float>} where ``expected`` is
computed by a pure-Python reference implementation identical to the Go
``ledger.Percentile`` (linear interpolation between closest ranks, i.e.
NumPy's default method="linear", C=1).

Re-running with the same seed produces byte-identical JSON, so the committed
data/pct_cases.json is exactly what this script emits.
"""

import json
import math
import os
import random

SEED = 20240804
NUM_RANDOM_CASES = 60


def percentile(xs, p):
    """Reference implementation mirroring ledger.Percentile in Go."""
    n = len(xs)
    if n == 0:
        return 0.0

    s = sorted(xs)
    if n == 1:
        return float(s[0])

    # Clamp p to the [0, 100] domain.
    if p < 0:
        p = 0.0
    elif p > 100:
        p = 100.0

    rank = (p / 100.0) * (n - 1)
    lo = math.floor(rank)
    hi = math.ceil(rank)
    frac = rank - lo
    return float(s[lo]) + frac * float(s[hi] - s[lo])


def build_cases():
    rng = random.Random(SEED)
    cases = []

    # A few deterministic hand-picked cases covering the boundaries.
    boundary_slices = [
        [5],
        [1, 2, 3, 4],
        [4, 3, 2, 1],
        [10, 20, 30],
    ]
    for xs in boundary_slices:
        for p in (0.0, 50.0, 100.0):
            cases.append({"xs": list(xs), "p": p, "expected": percentile(xs, p)})

    # Seeded random cases with varying lengths, values and p.
    for _ in range(NUM_RANDOM_CASES):
        n = rng.randint(1, 25)
        xs = [rng.randint(-1000, 1000) for _ in range(n)]
        # Mix boundary and random p values.
        choice = rng.random()
        if choice < 0.15:
            p = 0.0
        elif choice < 0.30:
            p = 100.0
        elif choice < 0.45:
            p = 50.0
        else:
            p = round(rng.uniform(0.0, 100.0), 6)
        cases.append({"xs": xs, "p": p, "expected": percentile(xs, p)})

    return cases


def main():
    repo_root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    out_dir = os.path.join(repo_root, "data")
    os.makedirs(out_dir, exist_ok=True)
    out_path = os.path.join(out_dir, "pct_cases.json")

    cases = build_cases()

    with open(out_path, "w", encoding="utf-8") as f:
        json.dump(cases, f, indent=2, sort_keys=True)
        f.write("\n")

    print(f"wrote {len(cases)} cases to {out_path}")


if __name__ == "__main__":
    main()
