# Performance

This document describes GitForge's performance targets, benchmark methodology, and measured results.

## Target Metrics

These are the performance targets GitForge aims to achieve. They are **goals, not claims** — actual performance depends on hardware, repository size, and workload.

| Operation | Target | Repository Size | Notes |
|---|---|---|---|
| `status` | < 100ms | 10,000 files | Cold (no cache) |
| `status` | < 20ms | 10,000 files | Warm (cached) |
| `log` (100 commits) | < 50ms | Any | Linear traversal |
| `log` (10,000 commits) | < 500ms | Any | With graph |
| `diff` (single file) | < 10ms | < 1MB file | Myers algorithm |
| `diff` (working tree) | < 200ms | 10,000 files | Full tree walk |
| `commit` | < 50ms | 100 staged files | Write tree + commit object |
| `add` (single file) | < 5ms | Any file < 10MB | Hash + index update |
| `init` | < 10ms | N/A | Create .git structure |
| Binary startup | < 20ms | N/A | To first useful output |
| Memory (idle TUI) | < 30MB | 10,000 files | Resident set size |
| Binary size | < 20MB | N/A | Stripped, single platform |

## Design for Performance

### Principles

1. **Minimize allocations.** Use sync.Pool, pre-allocated slices, and buffer reuse where it matters.
2. **Lazy loading.** Don't read objects until they're needed. Don't parse what you don't use.
3. **Streaming I/O.** Process large files (pack files, diffs) via streaming rather than loading into memory.
4. **Concurrency where beneficial.** Use goroutines for I/O-bound operations (tree walking, remote fetches), but avoid unnecessary parallelism for CPU-bound tasks on small inputs.
5. **Cache intelligently.** Cache parsed objects, pack file indexes, and directory listings. Invalidate on writes.
6. **Avoid CGO.** Pure Go eliminates CGO overhead and enables escape analysis and inlining optimizations across the entire call stack.

### Hot Paths

These code paths are performance-critical and receive extra attention:

- **Tree walking** — Recursive directory traversal for `status` and `diff`. Optimized with `os.ReadDir` and parallel stat calls.
- **Object hashing** — SHA-1 computation for every object. Uses `crypto/sha1` with buffer pooling.
- **Index parsing** — Binary index file parsing. Uses `encoding/binary` with pre-allocated buffers.
- **Diff computation** — Myers diff algorithm. Optimized for the common case of small changes in large files.
- **Pack file lookup** — Fanout table and binary search in pack index files.

## Benchmark Methodology

### Environment

- Benchmarks run on CI (GitHub Actions `ubuntu-latest`) for consistency.
- Local benchmarks should be run on a quiet machine with no background load.
- Each benchmark runs a minimum of 3 iterations; reported values are the median.

### Benchmark Suite

Benchmarks are located in `*_test.go` files alongside the code they measure, using Go's built-in `testing.B` framework.

```bash
# Run all benchmarks
CGO_ENABLED=0 go test -bench=. -benchmem ./...

# Run benchmarks for a specific package
CGO_ENABLED=0 go test -bench=. -benchmem ./internal/core/diff/

# Run benchmarks with CPU profiling
CGO_ENABLED=0 go test -bench=. -cpuprofile=cpu.prof ./internal/core/object/
```

### Profiling

```bash
# CPU profile
go tool pprof cpu.prof

# Memory profile
CGO_ENABLED=0 go test -bench=. -memprofile=mem.prof ./internal/core/object/
go tool pprof mem.prof

# Trace
CGO_ENABLED=0 go test -bench=. -trace=trace.out ./internal/core/object/
go tool trace trace.out
```

## Measured Results

> **No measurements yet.** This section will be populated once the core implementation is complete and benchmarks are written.

### Placeholder

| Operation | Measured | Target | Status |
|---|---|---|---|
| `status` (10k files, cold) | — | < 100ms | ⏳ Not measured |
| `status` (10k files, warm) | — | < 20ms | ⏳ Not measured |
| `log` (100 commits) | — | < 50ms | ⏳ Not measured |
| `diff` (single file) | — | < 10ms | ⏳ Not measured |
| Binary startup | — | < 20ms | ⏳ Not measured |
| Memory (idle TUI) | — | < 30MB | ⏳ Not measured |
| Binary size | — | < 20MB | ⏳ Not measured |

Results will be updated as benchmarks are implemented and run.
