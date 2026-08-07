#!/usr/bin/env python3
"""Generate data/pct_cases.json, a deterministic corpus of percentile cases.

Each case records an input slice, a percentile p, and the expected result
computed by a pure-Python reference implementation that mirrors the Go
ledger.Percentile function exactly (linear interpolation between the two
closest ranks, i.e. numpy's default "linear" method).

The generator seeds Python's random module with a fixed constant so that
re-running it emits byte-for-byte identical JSON. Run it from anywhere:

    python3 gen/gen_pct.py

and it rewrites data/pct_cases.json next to the repository root.
"""

import json
import math
import os
import random

SEED = 1234


def reference_percentile(xs, p):
    """Reference percentile matching ledger.Percentile.

    Empty input returns 0.0; p is clamped to [0, 100]; the input list is not
    mutated (a sorted copy is used).
    """
    if len(xs) == 0:
        return 0.0

    s = sorted(xs)
    if p < 0:
        p = 0.0
    elif p > 100:
        p = 100.0

    n = len(s)
    if n == 1:
        return float(s[0])

    r = (p / 100.0) * (n - 1)
    lo = math.floor(r)
    hi = math.ceil(r)
    frac = r - lo
    return float(s[lo]) + frac * (float(s[hi]) - float(s[lo]))


def build_cases():
    rng = random.Random(SEED)
    cases = []

    # Deterministic edge cases exercising the documented scenarios.
    cases.append({"xs": [], "p": 50.0})
    cases.append({"xs": [7], "p": 0.0})
    cases.append({"xs": [7], "p": 42.5})
    cases.append({"xs": [7], "p": 100.0})

    # A spread of randomised cases: varying slice lengths, value ranges, and
    # percentiles (including the 0 and 100 extremes).
    for _ in range(80):
        length = rng.randint(2, 40)
        hi = rng.choice([10, 100, 1000, 10000])
        lo = rng.choice([-1000, -100, 0])
        xs = [rng.randint(lo, hi) for _ in range(length)]
        p = rng.choice(
            [0.0, 100.0, rng.uniform(0.0, 100.0), rng.uniform(0.0, 100.0)]
        )
        cases.append({"xs": xs, "p": round(p, 6)})

    for case in cases:
        case["expected"] = reference_percentile(case["xs"], case["p"])

    return cases


def main():
    cases = build_cases()
    repo_root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    data_dir = os.path.join(repo_root, "data")
    os.makedirs(data_dir, exist_ok=True)
    out_path = os.path.join(data_dir, "pct_cases.json")
    with open(out_path, "w") as f:
        json.dump(cases, f, indent=2, sort_keys=True)
        f.write("\n")


if __name__ == "__main__":
    main()
