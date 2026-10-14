# GitForge Architecture

## Overview

GitForge is designed as a layered system where each layer has clear responsibilities and well-defined interfaces. The architecture prioritizes performance, testability, and extensibility.

## Layer Diagram

```
┌──────────────────────────────────────────────────────────────────┐
│                      Presentation Layer                         │
│                                                                  │
│  ┌──────────┐    ┌──────────────┐    ┌────────────────────┐     │
│  │   CLI    │    │     TUI      │    │     Web UI         │     │
│  │ (cobra)  │    │ (bubbletea)  │    │ (embedded SPA)     │     │
│  └────┬─────┘    └──────┬───────┘    └────────┬───────────┘     │
│       │                 │                      │                 │
├───────┴─────────────────┴──────────────────────┴─────────────────┤
│                      Application Layer                           │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐   │
│  │   Commands   │  │   Services   │  │     AI Engine        │   │
│  │              │  │              │  │ (commit msg, review)  │   │
│  └──────┬───────┘  └──────┬───────┘  └──────────┬───────────┘   │
│         │                 │                      │               │
├─────────┴─────────────────┴──────────────────────┴───────────────┤
│                        Core Layer                                │
│                                                                  │
│  ┌────────┐ ┌──────┐ ┌───────┐ ┌──────┐ ┌───────┐ ┌─────────┐ │
│  │Objects │ │ Refs │ │ Index │ │ Diff │ │ Merge │ │ Config  │ │
│  │(blob,  │ │      │ │       │ │      │ │       │ │         │ │
│  │ tree,  │ │      │ │       │ │      │ │       │ │         │ │
│  │commit, │ │      │ │       │ │      │ │       │ │         │ │
│  │ tag)   │ │      │ │       │ │      │ │       │ │         │ │
│  └────┬───┘ └──┬───┘ └───┬───┘ └──┬───┘ └───┬───┘ └────┬────┘ │
│       │        │         │        │         │           │       │
├───────┴────────┴─────────┴────────┴─────────┴───────────┴───────┤
│                      Storage Layer                               │
│                                                                  │
│  ┌───────────────────┐  ┌──────────────────┐  ┌──────────────┐  │
│  │   Object Store    │  │   Pack Files     │  │   Loose      │  │
│  │   (read/write)    │  │   (deltified)    │  │   Objects    │  │
│  └───────────────────┘  └──────────────────┘  └──────────────┘  │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

## Module Descriptions

### Presentation Layer

| Module | Package | Description |
|---|---|---|
| CLI | `cmd/gitforge`, `internal/cli` | Command-line interface built with Cobra. Parses arguments, invokes application-layer commands, and formats output. |
| TUI | `internal/tui` | Interactive terminal UI built with Bubble Tea. Provides a visual, keyboard-driven interface for all Git operations. |
| Web UI | `internal/web` | Embedded single-page application served via a local HTTP server. Communicates with the backend via a REST/WebSocket API. |

### Application Layer

| Module | Package | Description |
|---|---|---|
| Commands | `internal/app` | High-level command orchestration. Each user-facing operation (status, commit, push, etc.) is implemented as an application command that coordinates core-layer primitives. |
| Services | `internal/service` | Cross-cutting services: authentication, remote communication (HTTP/SSH), forge API clients. |
| AI Engine | `internal/ai` | Local-first AI assistant for commit message generation, code review suggestions, and PR summaries. Supports local GGUF models, HTTP APIs, and WASM runtimes. |

### Core Layer

| Module | Package | Description |
|---|---|---|
| Objects | `internal/core/object` | Git object model: blobs, trees, commits, and tags. Handles serialization, hashing, and verification. |
| Refs | `internal/core/ref` | Reference management: branches, tags, HEAD, symbolic refs. Implements reflog. |
| Index | `internal/core/index` | Staging area (index file) reading and writing. Tracks file state for commits. |
| Diff | `internal/core/diff` | Diff computation using Myers algorithm. Supports unified diff output and stat summaries. |
| Merge | `internal/core/merge` | Three-way merge with conflict detection and resolution helpers. |
| Config | `internal/config` | Configuration loading from TOML files and environment variables, with layered precedence. |

### Storage Layer

| Module | Package | Description |
|---|---|---|
| Object Store | `internal/store` | Unified interface for reading and writing Git objects, abstracting over loose objects and pack files. |
| Pack Files | `internal/store/pack` | Pack file reading, writing, and indexing. Supports deltified objects for efficient storage. |
| Loose Objects | `internal/store/loose` | Loose object storage — individual zlib-compressed files in `.git/objects`. |

## Data Flow

### Typical Command Flow

```
User Input → CLI Parser → Application Command → Core Operations → Storage I/O → Output
```

### Example: `gitforge status`

1. CLI parses `status` subcommand and flags.
2. Application layer's `StatusCommand` is invoked.
3. Core layer reads the HEAD ref, loads the current index, and walks the working tree.
4. Diff engine compares: HEAD tree vs. index (staged changes), index vs. working tree (unstaged changes).
5. Results are formatted and returned to the presentation layer.
6. CLI renders the output to the terminal.

## Design Decisions

| Decision | Rationale |
|---|---|
| Pure Go, no CGO | Enables easy cross-compilation, static binaries, and eliminates C dependency issues. |
| No libgit2 | Full control over the Git implementation; avoids CGO; educational value; tailor-made for our use cases. |
| Layered architecture | Clear separation of concerns; each layer is independently testable; presentation layer is swappable. |
| Interfaces for storage | Allows mock implementations for testing; future support for alternative backends. |
| TOML for configuration | Human-readable, well-specified, good Go library support, supports comments (unlike JSON). |
| Cobra for CLI | Industry standard, excellent subcommand support, auto-generated help and completions. |
| Bubble Tea for TUI | Elm-architecture model is clean and testable; excellent Go ecosystem support. |
| Local-first AI | Privacy-preserving; works offline; no API keys required for basic functionality. |
