#!/usr/bin/env python3
"""Deterministic generator for percentile test cases.

Writes data/pct_cases.json with a set of seeded random cases. Each case holds
an input slice ("xs"), a percentile position ("p", on the domain [0, 1]) and
the "expected" value produced by the reference implementation below.

The reference implementation mirrors ledger.Percentile in percentile.go
byte-for-byte in its numeric semantics: linear interpolation between closest
ranks (numpy "linear" / method C=1), safe on empty input, clamping p to
[0, 1], and operating on a sorted copy so the input is never mutated.

Run it from the repository root:

    python3 gen/gen_pct.py

Re-running always emits identical bytes because the RNG is seeded with a
fixed constant.
"""

import json
import math
import os
import random

# Fixed seed so the generated JSON is fully deterministic and re-runnable.
SEED = 20240521


def percentile(xs, p):
    """Reference percentile matching ledger.Percentile exactly."""
    n = len(xs)
    if n == 0:
        return 0.0

    sorted_xs = sorted(xs)

    if n == 1:
        return float(sorted_xs[0])

    # Clamp p to [0, 1].
    if p < 0:
        p = 0.0
    elif p > 1:
        p = 1.0

    r = p * (n - 1)
    lo = int(math.floor(r))
    hi = int(math.ceil(r))
    frac = r - lo

    if lo == hi:
        return float(sorted_xs[lo])
    return float(sorted_xs[lo]) + frac * float(sorted_xs[hi] - sorted_xs[lo])


def build_cases():
    rng = random.Random(SEED)
    cases = []

    # Fixed boundary/edge cases exercising the documented contract.
    cases.append({"xs": [], "p": 0.5})
    cases.append({"xs": [7], "p": 0.0})
    cases.append({"xs": [7], "p": 0.5})
    cases.append({"xs": [7], "p": 1.0})
    cases.append({"xs": [-3], "p": 0.25})

    # A grid of p values used across many random slices, including the
    # boundaries 0 and 1 and several interior points.
    p_values = [0.0, 0.1, 0.25, 0.33, 0.5, 0.66, 0.75, 0.9, 1.0]

    # Random slices spanning varied lengths, value ranges (with negatives),
    # and duplicates, each paired with several p values.
    lengths = [2, 3, 4, 5, 8, 13, 21, 50, 100]
    for length in lengths:
        for p in p_values:
            xs = [rng.randint(-100, 100) for _ in range(length)]
            cases.append({"xs": xs, "p": round(p, 4)})

    # Slices with heavy duplication to exercise the interpolation collapse.
    for _ in range(8):
        length = rng.randint(4, 20)
        pool = [rng.randint(-10, 10) for _ in range(3)]
        xs = [rng.choice(pool) for _ in range(length)]
        p = round(rng.random(), 4)
        cases.append({"xs": xs, "p": p})

    # A handful of fully random p values on random slices.
    for _ in range(10):
        length = rng.randint(2, 40)
        xs = [rng.randint(-1000, 1000) for _ in range(length)]
        p = round(rng.random(), 4)
        cases.append({"xs": xs, "p": p})

    for case in cases:
        case["expected"] = percentile(case["xs"], case["p"])

    return cases


def main():
    cases = build_cases()

    repo_root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    data_dir = os.path.join(repo_root, "data")
    os.makedirs(data_dir, exist_ok=True)
    out_path = os.path.join(data_dir, "pct_cases.json")

    with open(out_path, "w", encoding="utf-8") as f:
        json.dump(cases, f, indent=2, sort_keys=True)
        f.write("\n")

    print(f"wrote {len(cases)} cases to {out_path}")


if __name__ == "__main__":
    main()
