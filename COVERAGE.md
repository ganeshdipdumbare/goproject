# Test coverage

This document describes how test coverage is measured for `goproject`. Coverage
is produced entirely with the standard Go toolchain — `go test` and
`go tool cover` — so there is nothing extra to install and no third-party
coverage service involved. The figures in [Current coverage](#current-coverage)
are a point-in-time snapshot and will drift as code and tests are added; re-run
the commands below to get the current picture.

## Generating a coverage profile

Write a coverage profile for every package in the module:

```sh
go test ./... -coverprofile=coverage.out
```

If you only want the per-package percentages and no profile file:

```sh
go test ./... -cover
```

## Reading the results

Per-function coverage, with a trailing `total:` line for the whole profile:

```sh
go tool cover -func=coverage.out
```

An annotated, browsable HTML report showing which lines were and were not
executed:

```sh
go tool cover -html=coverage.out -o coverage.html
```

Omitting `-o` makes `go tool cover` open the report in your browser directly
instead of writing a file.

## Useful flags

- `-covermode=atomic` — use this when combining coverage with `-race`, so
  counter updates are safe under concurrent access:
  `go test ./... -race -covermode=atomic -coverprofile=coverage.out`.
- `-coverpkg=./...` — attribute coverage to packages other than the one under
  test, which is what you want for integration-style tests that exercise code
  across package boundaries.
- Write the profile outside the working tree, e.g.
  `go test ./... -coverprofile=/tmp/coverage.out`, to avoid leaving stray
  artifacts in the checkout.

## Current coverage

Captured on 2026-09-07 with `go version go1.27.0 linux/amd64` (the module
declares `go 1.22` in `go.mod`).

`go test ./... -coverprofile=coverage.out`:

```
	goproject		coverage: 0.0% of statements
```

`go tool cover -func=coverage.out`:

```
goproject/store.go:14:	NewStore	0.0%
goproject/store.go:21:	Get		0.0%
goproject/store.go:29:	Set		0.0%
goproject/store.go:36:	Keys		0.0%
total:			(statements)	0.0%
```

| Package (import path) | Statement coverage | Test files |
| --- | --- | --- |
| `goproject` | 0.0% | none |

The module currently contains a single package — `store.go`, declared as
`package ledger` and reachable at import path `goproject` — and it has no
`_test.go` files, so none of its statements are exercised and coverage is
0.0%. Once tests are added, re-run the commands above and update this section
with the real output.

## Generated artifacts

`coverage.out` and `coverage.html` are generated files and should not be
committed. The simplest habit is to write the profile to a temporary path, as
shown under [Useful flags](#useful-flags), so nothing lands in the working tree
in the first place; otherwise delete the files once you are done with them.

This is deliberately advice rather than an enforced rule: adding or changing
repository ignore rules is intentionally out of scope for this document.
