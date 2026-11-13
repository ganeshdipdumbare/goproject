# Benchmark and profiling helpers for the ledger store.
#
# benchstat (golang.org/x/perf/cmd/benchstat) is run via `go run` at a pinned
# version so it is not added to go.mod. It requires a Go 1.26+ toolchain; with
# the default GOTOOLCHAIN=auto, Go downloads one automatically if needed.

COUNT ?= 10
BENCH ?= .
BENCHSTAT_VERSION ?= v0.0.0-20261009192801-be2c69fb417e

BENCH_DIR := bench
BENCHSTAT := go run golang.org/x/perf/cmd/benchstat@$(BENCHSTAT_VERSION)

.PHONY: bench bench-baseline bench-stat bench-compare profile clean-bench

## bench: run all benchmarks COUNT times and save results to bench/new.txt
bench:
	@mkdir -p $(BENCH_DIR)
	@# Stream output via a temp file and only replace new.txt if go test
	@# succeeds, so a failed run neither reports success nor clobbers results.
	@{ go test -run='^$$' -bench='$(BENCH)' -benchmem -count=$(COUNT) ./... ; \
		echo $$? > $(BENCH_DIR)/.status; } | tee $(BENCH_DIR)/new.txt.tmp; \
		status=$$(cat $(BENCH_DIR)/.status); rm -f $(BENCH_DIR)/.status; \
		if [ "$$status" -ne 0 ]; then rm -f $(BENCH_DIR)/new.txt.tmp; \
			echo "go test failed (exit $$status); $(BENCH_DIR)/new.txt left unchanged" >&2; exit "$$status"; fi; \
		mv $(BENCH_DIR)/new.txt.tmp $(BENCH_DIR)/new.txt

## bench-baseline: save the latest results (bench/new.txt) as bench/old.txt
bench-baseline:
	@test -f $(BENCH_DIR)/new.txt || { echo "$(BENCH_DIR)/new.txt not found: run 'make bench' first" >&2; exit 1; }
	cp $(BENCH_DIR)/new.txt $(BENCH_DIR)/old.txt

## bench-stat: summarise bench/new.txt with benchstat
bench-stat:
	@test -f $(BENCH_DIR)/new.txt || { echo "$(BENCH_DIR)/new.txt not found: run 'make bench' first" >&2; exit 1; }
	$(BENCHSTAT) $(BENCH_DIR)/new.txt

## bench-compare: compare bench/old.txt against bench/new.txt with benchstat
bench-compare:
	@test -f $(BENCH_DIR)/old.txt || { echo "$(BENCH_DIR)/old.txt not found: run 'make bench' then 'make bench-baseline' first" >&2; exit 1; }
	@test -f $(BENCH_DIR)/new.txt || { echo "$(BENCH_DIR)/new.txt not found: run 'make bench' first" >&2; exit 1; }
	$(BENCHSTAT) $(BENCH_DIR)/old.txt $(BENCH_DIR)/new.txt

## profile: write CPU and memory profiles to bench/ for go tool pprof
profile:
	@mkdir -p $(BENCH_DIR)
	go test -run='^$$' -bench='$(BENCH)' -benchmem -cpuprofile $(BENCH_DIR)/cpu.out -memprofile $(BENCH_DIR)/mem.out -o $(BENCH_DIR)/ledger.test .
	@echo "Inspect with: go tool pprof -top $(BENCH_DIR)/cpu.out"
	@echo "          or: go tool pprof -top $(BENCH_DIR)/mem.out"

## clean-bench: remove benchmark results and profiles
clean-bench:
	rm -rf $(BENCH_DIR)
