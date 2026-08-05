#!/usr/bin/env python3
"""Generate a deterministic corpus of percentile test cases.

Writes data/pct_cases.json relative to the repository root. The reference
implementation mirrors ledger.Percentile in percentile.go exactly (linear
interpolation between the two closest ranks, the numpy "linear"/type-7
method) so the Go tests and this generator agree to floating-point epsilon
without any external dependencies.

Re-running this script produces byte-identical output.
"""

import json
import math
import os
import random

SEED = 20240517


def reference_percentile(xs, p):
    """Pure-Python mirror of ledger.Percentile."""
    n = len(xs)
    if n == 0:
        return 0.0

    s = sorted(xs)
    if n == 1:
        return float(s[0])

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

    # Explicit edge cases.
    cases.append({"xs": [42], "p": 50.0})
    cases.append({"xs": [7, 7, 7, 7], "p": 25.0})
    cases.append({"xs": [-5, -1, 0, 3, 9], "p": 0.0})
    cases.append({"xs": [-5, -1, 0, 3, 9], "p": 100.0})

    # Random cases spanning varied lengths, value ranges, and p values.
    boundary_ps = [0.0, 25.0, 50.0, 75.0, 100.0]
    for i in range(70):
        length = rng.randint(1, 20)
        lo = rng.randint(-100, 0)
        hi = rng.randint(1, 100)
        xs = [rng.randint(lo, hi) for _ in range(length)]

        if i % 5 == 0:
            p = boundary_ps[(i // 5) % len(boundary_ps)]
        else:
            p = round(rng.uniform(0.0, 100.0), 4)

        cases.append({"xs": xs, "p": p})

    for case in cases:
        case["expected"] = reference_percentile(case["xs"], case["p"])

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
