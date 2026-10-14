# Dependencies

This document explains the rationale behind each dependency used in GitForge.

## Philosophy

GitForge aims for a **minimal dependency footprint**. Every dependency must justify its inclusion by providing significant value that would be costly or error-prone to implement from scratch. We prefer the standard library where feasible.

## Direct Dependencies

> **Note:** Dependencies will be added incrementally as features are implemented. This document will be updated with each addition.

### Planned Dependencies

| Dependency | Purpose | Justification |
|---|---|---|
| `github.com/spf13/cobra` | CLI framework | Industry-standard CLI library. Provides subcommand routing, flag parsing, auto-generated help, shell completions. Reimplementing this well would take significant effort with little benefit. |
| `github.com/BurntSushi/toml` | TOML configuration parsing | Well-tested, spec-compliant TOML parser. Configuration parsing is a solved problem. |
| `github.com/charmbracelet/bubbletea` | TUI framework | Elm-architecture TUI framework. The best-in-class Go TUI library. Building a comparable framework from scratch would be a project in itself. |
| `github.com/charmbracelet/lipgloss` | TUI styling | Companion to Bubble Tea for terminal styling. Handles ANSI escape sequences, color profiles, and layout. |
| `github.com/charmbracelet/bubbles` | TUI components | Pre-built Bubble Tea components (text input, viewport, list, etc.). Avoids reinventing common UI widgets. |

### Development Dependencies

| Dependency | Purpose | Justification |
|---|---|---|
| `golangci-lint` | Linting | Aggregates dozens of Go linters into a single tool. Essential for code quality. |
| `goreleaser` | Release automation | Handles cross-compilation, packaging, checksums, and changelog generation for releases. |

## Dependency Policy

1. **Evaluate before adding.** Every new dependency must be discussed and justified.
2. **Prefer the standard library.** If the standard library can do it reasonably well, use it.
3. **Audit transitively.** Consider the full transitive dependency tree, not just the direct import.
4. **Pin versions.** All dependencies are pinned via `go.sum`.
5. **No CGO.** Dependencies must not require CGO. This ensures easy cross-compilation and static binaries.
6. **License compatibility.** All dependencies must use MIT, BSD, Apache 2.0, or similarly permissive licenses.
7. **Active maintenance.** Prefer dependencies that are actively maintained and widely used.

## Why Not libgit2?

GitForge implements Git operations in pure Go rather than binding to libgit2 via CGO. Reasons:

- **Cross-compilation**: CGO breaks `GOOS`/`GOARCH` cross-compilation. Pure Go enables trivial multi-platform builds.
- **Static binaries**: No shared library dependencies at runtime.
- **Full control**: We can optimize for our specific use cases and add features without waiting for upstream.
- **Learning**: Building Git internals from scratch is a core goal of this project.
- **Debugging**: Pure Go stack traces are easier to debug than mixed Go/C stacks.
