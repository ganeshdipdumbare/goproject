# Test coverage

## Overview

This document describes the test-coverage story for this repository: how to
produce a coverage profile locally, how to read it, and what the current
numbers are. Figures below were measured against commit `999ebba` (2026-07-22)
with the Go toolchain declared in `go.mod` (Go 1.22 or newer).

## Generating a coverage report

Produce a coverage profile for every package in the module:

```bash
go test ./... -coverprofile=coverage.out -covermode=atomic
```

`-covermode=atomic` is used because the code under test (`Store`) is designed
for concurrent use and the tests may exercise it from multiple goroutines;
atomic counters keep the profile accurate under `-race` and parallel tests.

Read the profile as a per-function summary:

```bash
go tool cover -func=coverage.out
```

Or open an annotated, line-by-line HTML view in your browser:

```bash
go tool cover -html=coverage.out
```

To write the HTML report to a file instead of opening a browser:

```bash
go tool cover -html=coverage.out -o coverage.html
```

## Current coverage

| Package | Statements covered | Notes |
| --- | --- | --- |
| `goproject` (root, package `ledger`) | 0.0% | No test files yet; `store.go` holds `NewStore`, `Get`, `Set` and `Keys`. |
| **Total** | **0.0%** | Single-package module. |

Per-function breakdown from `go tool cover -func=coverage.out`:

| Function | File | Statements covered |
| --- | --- | --- |
| `NewStore` | `store.go:14` | 0.0% |
| `Get` | `store.go:21` | 0.0% |
| `Set` | `store.go:29` | 0.0% |
| `Keys` | `store.go:36` | 0.0% |

`go test ./...` currently reports `?   goproject   [no test files]`, so the
0.0% figures reflect an absence of tests rather than failing or skipped ones.

## Interpreting the numbers

Statement coverage is a signal, not a goal: 100% coverage does not prove the
behaviour is correct, and a lower number does not automatically mean the code
is untrustworthy. It is most useful for spotting code paths that no test ever
executes.

For this repository there are no deliberately excluded packages — there is no
generated code, no `main` wiring and no mocks. Every uncovered statement listed
above is ordinary library code that is simply not exercised yet. The obvious
first tests are: a round-trip `Set`/`Get`, a `Get` on a missing key returning
`(0, false)`, and `Keys` returning a snapshot of the stored keys.

## Housekeeping

`coverage.out` and `coverage.html` are build artifacts and must not be
committed. This repository **intentionally does not add ignore rules for
them** (see issue #27), so they will show up as untracked files in
`git status`. Clean them up manually before committing:

```bash
rm -f coverage.out coverage.html
```

## Refreshing this document

Regenerate the numbers and update the tables above whenever coverage moves
materially (a new package, new tests, or a noticeable swing in a percentage):

```bash
go test ./... -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out
```

Remember to update the "as of" commit and date in the overview, then delete the
generated artifacts.
