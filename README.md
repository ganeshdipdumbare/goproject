# goproject

## Performance

Every exported function in `store.go` (`NewStore`, `Get`, `Set`, `Keys`) has a
benchmark in `store_bench_test.go`. Parallel variants (`BenchmarkGetParallel`,
`BenchmarkSetParallel`, `BenchmarkMixedParallel`) measure lock contention under
concurrent use.

Results are analysed with
[benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat), run through
`go run` at the version pinned in the `Makefile` (`BENCHSTAT_VERSION`). It is
not added to `go.mod`. benchstat needs a Go 1.26+ toolchain; with the default
`GOTOOLCHAIN=auto`, Go downloads one automatically if yours is older.

### Run the benchmarks

```sh
make bench               # runs all benchmarks COUNT times (default 10), saves bench/new.txt
make bench COUNT=3       # fewer samples for a quicker run
make bench BENCH=Get     # only benchmarks matching a regex
make bench-stat          # summarise bench/new.txt: sec/op, B/op, allocs/op with confidence intervals
```

### Compare two runs

```sh
make bench               # measure the current code
make bench-baseline      # save bench/new.txt as bench/old.txt
# ...make your change...
make bench               # measure again
make bench-compare       # benchstat bench/old.txt bench/new.txt: delta and p-value per benchmark
```

Results vary between machines and runs. Compare runs made on the same machine,
with `COUNT` of at least 10, and keep the machine otherwise idle.

### Profile

```sh
make profile                       # writes bench/cpu.out and bench/mem.out
go tool pprof -top bench/cpu.out   # top CPU consumers
go tool pprof -top bench/mem.out   # top allocations
```

`make clean-bench` removes the `bench/` directory.
