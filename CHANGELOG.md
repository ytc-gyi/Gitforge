# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Phase 2: Complete CLI**:
  - Implemented all remaining Git commands: `init`, `clone`, `branch` (`-d`, `-a`, `-r`), `checkout` (`-b`), `add` (`-p` with interactive hunk-level staging), `commit` (`-m`, `--amend`), `remote` (`add`, `remove`), `fetch` (`--all`, `--prune`, `--tags`), `pull` (`--rebase`), `push` (`-f`, `-u`, `--tags`), `stash` (`list`, `pop`, `apply`, `drop`), and `tag` (`-d`, `-a`, `-m`).
  - Hunk-level staging (`add -p`) with standalone unified diff patch generator and `git apply --cached`.
  - Shell completion generation (`gitforge completion [bash|zsh|fish|powershell]`).
  - Lip Gloss styled terminal output with automatic plain fallback when stdout is not a TTY or `NO_COLOR` is set.
  - Verified cross-compilation with `CGO_ENABLED=0` for all 9 target platforms (`linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`, `windows/arm64`, `freebsd/amd64`, `openbsd/amd64`, `android/arm64`).
  - Native Termux installation documentation and script (`scripts/install-termux.sh`).
  - Checksum-verified installation script (`scripts/install.sh`).
- **Phase 1: Foundation**:
  - Core domain types, sentinel errors, typed `GitError` and `ConfigError`.
  - Pub/sub `EventBus` and `ModuleRegistry` wired into `core.Engine` facade.
  - Configuration system in `internal/config` supporting TOML loading, precedence (env > local > user > system > defaults), env var expansion (`${VAR}`), typed validation, and secret token redaction.
  - Git backend abstraction (`git.Backend`) and machine-readable `CLIBackend`.
  - Robust machine-readable parsers for `status --porcelain=v2 -z`, `log`, `diff` (with hunk and diffstat parsing), and `branch`.
  - Pure commit graph layout algorithm with lane assignment in `internal/git/graph`.
  - CLI commands via Cobra: `version`, `doctor`, `status` (`--json`), `log` (`--oneline`, `--graph`, `--json`), `diff`, `show`, `config` (get/set/list with redaction), and transparent Git passthrough.
  - Comprehensive unit, table-driven, fuzz (`FuzzParseStatus`, `FuzzParseLog`, `FuzzParseDiff`), and real Git integration tests with race detector enabled.
  - Full repo scaffolding: Makefile, `.golangci.yml`, `.goreleaser.yaml`, Dockerfile, GitHub Actions workflows (`ci.yml`, `release.yml`, `codeql.yml`), docs and licensing.
