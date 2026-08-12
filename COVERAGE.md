# Test Coverage

This document explains how test coverage works for the `goproject` repository,
how to generate a coverage report with the standard Go toolchain, and how to
read the results.

## Generating a coverage profile

Run the test suite with the `-coverprofile` flag to write a coverage profile to
disk:

```sh
go test ./... -coverprofile=coverage.out
```

This runs every package's tests and records which statements were executed into
`coverage.out`.

## Viewing a summary

To print a per-function and total coverage summary to the terminal:

```sh
go tool cover -func=coverage.out
```

The final `total:` line reports the overall statement coverage for the module.

## Viewing an annotated HTML report

To open an annotated, line-by-line HTML report that highlights covered and
uncovered code:

```sh
go tool cover -html=coverage.out
```

To write the report to a file instead of opening a browser:

```sh
go tool cover -html=coverage.out -o coverage.html
```

## Interpreting the results

- **Covered lines** are those executed by at least one test.
- **Uncovered lines** are shown in red in the HTML report — these are code
  paths no test currently exercises.
- The percentage is *statement* coverage: the fraction of statements executed,
  not a guarantee that every input or edge case is tested.

## Current status

The repository currently ships the `Store` type in `store.go` but does not yet
include test files, so `go test ./... -cover` reports:

```
coverage: 0.0% of statements
```

Adding tests for `Store` (its `Get`, `Set`, and `Keys` methods) will raise this
figure. Re-run the commands above after adding tests to measure the new
coverage.
