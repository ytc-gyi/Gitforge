# GitForge Roadmap

## Current Status: Phase 2 Complete ✅ | Starting Phase 3 (TUI) 🚀

---

## Phase 1 — Foundation

**Status:** ✅ Completed

- [x] Repository scaffold, `go.mod` (Go 1.23+), multi-target Makefile, `.golangci.yml`, `.goreleaser.yaml`, Dockerfile.
- [x] CI workflows (Linux/macOS/Windows matrix, release, CodeQL).
- [x] Core domain types, errors (`GitError`, `ConfigError`), `EventBus` pub/sub, `ModuleRegistry`, and `core.Engine` facade.
- [x] `internal/config` with full precedence hierarchy (env > local > user > system > defaults), `${VAR}` expansion, typed validation, and secret redaction.
- [x] `internal/git` with `Backend` interface and `CLIBackend` implementation.
- [x] Machine-readable parsers for `status --porcelain=v2 -z`, `log`, `diff` (with hunk and diffstat parsing), and `branch`.
- [x] Pure commit graph layout algorithm with lane assignment in `internal/git/graph`.
- [x] CLI commands via Cobra: `version`, `doctor`, `status` (`--json`), `log` (`--oneline`, `--graph`, `--json`), `diff`, `show`, `config` (get/set/list with redaction), and unknown command Git passthrough.
- [x] Test suite: table-driven parser tests, fuzz tests (`FuzzParseStatus`, `FuzzParseLog`, `FuzzParseDiff`), and real temporary Git repository integration tests.
- [x] `make check` passing with race detector enabled; coverage ≥ 80% on config and git parsers.

---

## Phase 2 — Complete CLI

**Status:** ✅ Completed

- [x] All remaining Git commands: `init`, `clone`, `branch` (`-d`, `-a`, `-r`), `checkout` (`-b`), `add` (`-p` with hunk-level staging via generated patches), `commit` (`-m`, `--amend`), `remote` (`add`, `remove`), `fetch` (`--all`, `--prune`, `--tags`), `pull` (`--rebase`), `push` (`-f`, `-u`, `--tags`), `stash` (`list`, `pop`, `apply`, `drop`), `tag` (`-d`, `-a`, `-m`).
- [x] Interactive patch application (`add -p`) with individual hunk staging/discarding.
- [x] Full alias expansion and shell completions (bash, zsh, fish, powershell).
- [x] Enhanced terminal output with Lip Gloss and plain output when stdout is not a TTY or `NO_COLOR` is set.
- [x] Cross-compilation verified with `CGO_ENABLED=0` for linux/amd64+arm64, darwin/amd64+arm64, windows/amd64+arm64, freebsd/amd64, openbsd/amd64, android/arm64.
- [x] Native Termux install script (`scripts/install-termux.sh`) and checksum-verified installer (`scripts/install.sh`).

---

## Phase 3 — TUI (Terminal User Interface)

**Status:** 📋 Planned

- [ ] Multi-panel layout using Bubble Tea + Lip Gloss + Bubbles: header, left sidebar (Repo, Branches, Remotes, Actions), center (commit graph/log or status), right (diff).
- [ ] Keyboard navigation and keybindings table (Tab, j/k, Enter, Space, s, c, p, P, r, f, d, l, b, ?, q) with configurable keymaps.
- [ ] Mouse support: click to focus, wheel scroll, drag to resize.
- [ ] Interactive hunk staging and side-by-side conflict resolver ("accept ours / theirs / both / edit").
- [ ] Built-in theme system (Catppuccin, Dracula, Nord, Tokyo Night, Gruvbox, Everforest, Rose Pine, auto).
- [ ] Live updates using `fsnotify` debounced watcher.
- [ ] Model tests via `teatest` and golden files at 80x24 and 120x40.

---

## Phase 4 — Web UI

**Status:** 📋 Planned

- [ ] REST + SSE/WebSocket API mirroring engine queries and events.
- [ ] Embedded React 18 + TypeScript + Vite + Tailwind frontend served on localhost with random free port and session token.
- [ ] Feature parity with TUI: hunk staging, commit graph (SVG), unified and side-by-side diff, conflict resolution.
- [ ] PWA manifest, service worker, offline shell, mobile-responsive layout.
- [ ] OpenAPI specification in `docs/openapi.yaml`.

---

## Phase 5 — Plugins (WASM)

**Status:** 📋 Planned

- [ ] Sandboxed WASM runtime using pure-Go `wazero`.
- [ ] Plugin manifest (`plugin.toml`), capabilities model (`git.read`, `git.write`, `net.http:<host>`, `fs.read:<path>`, etc.).
- [ ] Detached Ed25519 signature verification against trusted sources.
- [ ] Resource quotas: memory limits, execution timeout.
- [ ] Plugin SDK and two reference plugins (conventional commits validator, webhook notifier).
- [ ] CLI commands: `plugin list`, `install`, `remove`, `update`, `info`.

---

## Phase 6 — Local AI Assistant

**Status:** 📋 Planned

- [ ] `InferenceEngine` interface with deterministic rule-based fallback (no model required).
- [ ] Locally running HTTP inference engine (llama.cpp / Ollama on 127.0.0.1).
- [ ] WASM llama.cpp engine under wazero (best effort).
- [ ] Verified resumable model download manager with SHA-256 integrity check.
- [ ] CLI commands: `ai commit`, `ai explain`, `ai review`, `ai suggest-branch`, `ai models`, `ai download`.
- [ ] Token-budgeted context builder and embedded prompt templates.
- [ ] Privacy invariant test (no outbound network calls from AI engine except local endpoint and downloader).

---

## Phase 7 — Forge Integrations

**Status:** 📋 Planned

- [ ] `ForgeProvider` interface with GitHub and GitLab (including self-hosted) implementations.
- [ ] PR/MR list, view, checkout, create; Issue list, create; CI status summary.
- [ ] Token authentication from config/env with redaction.
- [ ] Mock-server testing suite.

---

## Phase 8 — Release Engineering & Polish

**Status:** 📋 Planned

- [ ] GoReleaser release archives, checksums, SBOM, signatures.
- [ ] Minimal non-root Dockerfile.
- [ ] Verified installation scripts (`install.sh`, `install-termux.sh`).
- [ ] Packaging templates (Homebrew, Scoop, Nix, AUR).
- [ ] CI security audits: `govulncheck`, `gosec`, `CodeQL`.
- [ ] Benchmarks and `docs/performance.md` with measured numbers.

---

## Future & Out of Scope for v1

- **etcd / Distributed Sync:** Evaluated and excluded for v1. GitForge prioritizes single static binary local-first architecture without external daemon dependencies. A `SyncProvider` interface will be defined with no-op implementation.
- **Real-Time Multi-User Collaboration:** Evaluated and deferred to post-v1. The Web API is structured so a real-time CRDT or OT collaboration transport can be layered on top in the future.
- **Handwritten SIMD:** Premature optimization. Replaced with `sync.Pool` and zero-allocation parsing routines guided by benchmark profiling.
