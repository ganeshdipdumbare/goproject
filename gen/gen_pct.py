#!/usr/bin/env python3
"""Generate deterministic test cases for ledger.Percentile.

This writes data/pct_cases.json with a fixed, seeded set of cases so that
re-running the generator yields byte-identical output. Each case records the
input slice, the percentile p, and the expected value computed by the
reference implementation below.

The reference algorithm mirrors ledger.Percentile exactly: values are sorted
ascending and the p-th percentile is the linear interpolation between the two
closest ranks (numpy's default "linear" method) on the 0-indexed sorted data.
"""

import json
import math
import os
import random


def percentile(xs, p):
    """Reference implementation matching percentile.go."""
    if not xs:
        return 0.0
    if p < 0:
        p = 0.0
    if p > 100:
        p = 100.0

    s = sorted(xs)
    if len(s) == 1:
        return float(s[0])

    r = (len(s) - 1) * (p / 100.0)
    lo = int(math.floor(r))
    hi = int(math.ceil(r))
    frac = r - lo
    return float(s[lo]) + frac * float(s[hi] - s[lo])


def main():
    rng = random.Random(12345)

    # A spread of "round" percentiles plus arbitrary ones exercised below.
    round_ps = [0.0, 25.0, 50.0, 75.0, 90.0, 95.0, 99.0, 100.0]

    cases = []

    # Single-element slices: any p returns that element.
    for _ in range(6):
        val = rng.randint(-100, 100)
        p = rng.choice(round_ps + [rng.uniform(0.0, 100.0)])
        cases.append({"xs": [val], "p": p})

    # Small and large slices with round percentiles.
    for _ in range(30):
        n = rng.randint(2, 40)
        xs = [rng.randint(-1000, 1000) for _ in range(n)]
        p = rng.choice(round_ps)
        cases.append({"xs": xs, "p": p})

    # Slices with arbitrary fractional percentiles.
    for _ in range(30):
        n = rng.randint(2, 60)
        xs = [rng.randint(-1000, 1000) for _ in range(n)]
        p = rng.uniform(0.0, 100.0)
        cases.append({"xs": xs, "p": p})

    for c in cases:
        c["expected"] = percentile(c["xs"], c["p"])

    here = os.path.dirname(os.path.abspath(__file__))
    out_dir = os.path.join(here, "..", "data")
    os.makedirs(out_dir, exist_ok=True)
    out_path = os.path.join(out_dir, "pct_cases.json")

    with open(out_path, "w") as f:
        json.dump(cases, f, indent=2, sort_keys=True)
        f.write("\n")


if __name__ == "__main__":
    main()
