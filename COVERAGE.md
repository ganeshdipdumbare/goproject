# Test Coverage

This document explains how test coverage is measured for this repository and
records the current coverage snapshot for the Go codebase.

## Overview

`goproject` is a Go module. Coverage is measured with Go's built-in tooling
(`go test` and `go tool cover`); no third-party coverage tool is required.

## Generating a coverage report

Run the test suite with a coverage profile from the repository root:

```sh
go test ./... -coverprofile=coverage.out
```

This produces a `coverage.out` profile file describing which statements were
executed by the tests.

## Viewing the report

Print a per-function and total summary to the terminal:

```sh
go tool cover -func=coverage.out
```

Or open an interactive HTML report that highlights covered and uncovered lines:

```sh
go tool cover -html=coverage.out
```

To write the HTML report to a file instead of opening a browser:

```sh
go tool cover -html=coverage.out -o coverage.html
```

## Current coverage

Captured by running `go test ./... -coverprofile=coverage.out` followed by
`go tool cover -func=coverage.out`.

| Package     | Statement coverage |
| ----------- | ------------------ |
| `goproject` | 0.0%               |
| **Total**   | **0.0%**           |

Per-function breakdown for `goproject`:

| Function     | File            | Coverage |
| ------------ | --------------- | -------- |
| `NewStore`   | `store.go:14`   | 0.0%     |
| `Get`        | `store.go:21`   | 0.0%     |
| `Set`        | `store.go:29`   | 0.0%     |
| `Keys`       | `store.go:36`   | 0.0%     |

> The repository currently ships no `_test.go` files, so `go test ./...`
> reports `[no test files]` and coverage is `0.0%`. Adding tests for the
> `Store` type in `store.go` will raise these numbers; regenerate this report
> after adding tests to keep it accurate.

## A note on generated files

`coverage.out` and `coverage.html` are generated artifacts. Please do not
commit them to the repository — regenerate them locally with the commands
above whenever you need a fresh report.
