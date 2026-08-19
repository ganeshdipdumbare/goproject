# Test Coverage

This document explains how test coverage is measured for `goproject`, how to
generate and read coverage reports locally, and what the current coverage of
the repository looks like.

## Overview

`goproject` is a standard Go module (see `go.mod`), so coverage is measured
with the built-in Go toolchain — no third-party tools are required. Go's
`go test` command can instrument packages while running their tests and emit a
coverage profile that can then be summarised on the command line or rendered as
an annotated HTML report.

## Running coverage locally

To run the full test suite and write a coverage profile to `coverage.out`:

```sh
go test ./... -coverprofile=coverage.out
```

For a quick, per-package percentage without writing a profile file, use the
`-cover` flag on its own:

```sh
go test ./... -cover
```

## Reading the coverage summary

Once you have a `coverage.out` profile, print a per-function and total summary
with `go tool cover`:

```sh
go tool cover -func=coverage.out
```

This lists each function, its source location, and the percentage of its
statements exercised by the tests, followed by a `total` line for the whole
profile.

## Generating an HTML report

For a visual, line-by-line view of which statements are covered, render the
profile to HTML:

```sh
go tool cover -html=coverage.out -o coverage.html
```

Then open `coverage.html` in a browser. Covered lines are shown in green and
uncovered lines in red. Omit `-o coverage.html` to open the report directly in
your default browser instead of writing a file.

## Current coverage

The repository currently ships the `Store` implementation in `store.go` but
does **not** include any test files, so measured statement coverage is `0.0%`.
Running `go test ./... -coverprofile=coverage.out` followed by
`go tool cover -func=coverage.out` produces:

```text
goproject/store.go:14:  NewStore    0.0%
goproject/store.go:21:  Get         0.0%
goproject/store.go:29:  Set         0.0%
goproject/store.go:36:  Keys        0.0%
total:                  (statements)    0.0%
```

| Package     | Coverage |
| ----------- | -------- |
| `goproject` | 0.0%     |
| **Total**   | **0.0%** |

Add tests (for example a `store_test.go` exercising `NewStore`, `Get`, `Set`,
and `Keys`) and re-run the commands above to raise these numbers. Please update
this table whenever the coverage changes.

## Note on generated artifacts

The `coverage.out` profile and `coverage.html` report are produced locally when
you run the commands above. They are build artifacts rather than source and do
not need to be committed. This repository intentionally does not add a
`.gitignore` for them; exclude them through your own local or global Git
configuration if you prefer not to see them as untracked files.
