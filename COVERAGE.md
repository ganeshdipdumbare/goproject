# Test Coverage

This document is the single reference for `goproject`'s test-coverage
workflow: how to generate coverage locally with Go's built-in tooling, how to
read the output, and where the current coverage numbers live.

## Overview

`goproject` uses the standard Go toolchain (`go test` together with
`go tool cover`) for coverage — there is no third-party coverage runner to
install. Coverage is reported as **statement coverage**: the percentage of
executable statements that were run at least once by the test suite.

## Generating a coverage profile

From the repository root, run the full test suite and write a coverage profile:

```sh
go test ./... -coverprofile=coverage.out
```

This produces `coverage.out`, a machine-readable profile that the reporting
commands below consume. Re-run it whenever you want fresh numbers.

## Reading the summary report

For a per-function and total breakdown in your terminal:

```sh
go tool cover -func=coverage.out
```

Each line shows a function and its statement coverage; the final `total:` line
is the coverage across the whole module. That `total:` value is the single
number to quote when talking about overall coverage.

## Browsing the HTML report

For a line-by-line, colour-coded view of exactly which statements are covered:

```sh
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` in a browser. Covered statements are shown in green and
uncovered statements in red, making it easy to spot gaps. Omit `-o` to open the
report directly in your default browser instead of writing a file.

## Current coverage

The numbers below come from `go tool cover -func=coverage.out` against the
current `main`. Regenerate them with the commands above and update this table
when coverage changes.

| Package     | Statement coverage |
| ----------- | ------------------ |
| `goproject` | 0.0%               |
| **Total**   | **0.0%**           |

> The module currently ships no `_test.go` files, so measured coverage is
> `0.0%`. As tests are added, refresh the profile and update this table with
> the new `total:` figure.
