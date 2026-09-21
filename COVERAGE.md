# Test coverage

## Overview

This repository tracks statement coverage for its Go packages so that
contributors can see, at a glance, which code paths are exercised by the test
suite. Coverage is a useful signal, not a guarantee of correctness: a line can
be executed by a test without that test asserting anything meaningful about its
behaviour. Treat the numbers below as a way to find untested code, not as proof
that the tested code is right.

## Running coverage locally

Everything below only needs a working Go toolchain (the module targets Go 1.22).

Produce a coverage profile:

```sh
go test ./... -coverprofile=coverage.out -covermode=atomic
```

Print per-function and total coverage:

```sh
go tool cover -func=coverage.out
```

Render a browsable, line-by-line HTML report:

```sh
go tool cover -html=coverage.out -o coverage.html
```

There is no `Makefile` in this repository, so the commands above are the
canonical entry point.

## Current coverage

As of commit `999ebba` (2026-07-22):

| Package | Statements covered | Notes |
|---|---|---|
| `goproject` (package `ledger`, `store.go`) | 0.0% | No test files exist yet; `NewStore`, `Get`, `Set` and `Keys` are all uncovered. |
| **Total** | **0.0%** | Output of the `total:` line from `go tool cover -func=coverage.out`. |

Raw output of `go tool cover -func=coverage.out`:

```text
goproject/store.go:14:  NewStore        0.0%
goproject/store.go:21:  Get             0.0%
goproject/store.go:29:  Set             0.0%
goproject/store.go:36:  Keys            0.0%
total:                  (statements)    0.0%
```

This snapshot drifts as code and tests change — regenerate it with the commands
above before relying on it.

## Interpreting the numbers

- Go reports **statement** coverage, not branch or condition coverage. A 100%
  figure means every statement ran at least once, not that every input or error
  path was considered.
- `-covermode=atomic` is used so the counters stay correct when tests exercise
  code from multiple goroutines — relevant here, since `Store` is explicitly a
  concurrency-safe type and should be tested under `-race`.
- Wiring code (a future `main.go`), generated code and mocks typically pull the
  overall percentage down without saying much about test quality; judge those
  packages on their own rather than on the aggregate.
- The single package currently in the tree is at 0.0% purely because it has no
  tests yet, not because any code is intentionally excluded.

## Coverage goals

These are guidelines, not an enforced gate — nothing in this repository fails a
build on a coverage threshold today.

- Aim for **80%** statement coverage on packages containing business logic.
- New packages and new exported functions should ship with tests in the same
  change.
- Concurrency-sensitive code (such as `Store`) should additionally be covered by
  a test run with `go test -race ./...`.

## CI integration

There are no workflows under `.github/workflows/` in this repository, so
coverage is currently a **local-only** step. If CI is added later, the commands
in [Running coverage locally](#running-coverage-locally) are the ones it should
run, and this section should be updated to say where the report is published.

## Artifacts

`coverage.out` and `coverage.html` are build artifacts and must not be
committed. This repository deliberately has **no `.gitignore`**, so nothing
excludes them automatically. Either write them outside the working tree:

```sh
go test ./... -coverprofile="$TMPDIR/coverage.out" -covermode=atomic
```

or delete them once you are done:

```sh
rm -f coverage.out coverage.html
```

Check `git status` before committing to confirm neither file is staged.
