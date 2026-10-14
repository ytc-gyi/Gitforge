# GitForge

**A fast, cross-platform, local-first Git client that ships as a single static Go binary.**

GitForge offers three front ends over one shared, high-performance core:
- **CLI** (Cobra): Git-compatible command surface with styled output and seamless Git passthrough.
- **TUI** (Bubble Tea + Lip Gloss + Bubbles): Multi-panel terminal interface with vim-style keybindings.
- **Web UI** (React 18 + TypeScript + Vite + Tailwind): Embedded local web interface installable as an offline PWA.

On top of this sits a local-first AI assistant, a sandboxed WASM plugin system (wazero), and GitHub/GitLab forge integrations.

---

## Product Philosophy

- **Terminal-first, not terminal-only.** Work in your terminal, in an interactive TUI, or in a local browser UI with equal speed.
- **Local-first, cloud-optional.** Everything works offline. Nothing is sent anywhere without explicit opt-in.
- **AI that respects your code.** Inference runs locally. Zero telemetry, ever.
- **Extensible by default.** Features are modules behind clean interfaces; WASM plugins can extend them.
- **Runs everywhere.** Linux, macOS, Windows, WSL2, FreeBSD, OpenBSD, and Termux (Android). Single static binary with `CGO_ENABLED=0`.

---

## Features

| Interface / Module | Description | Engine |
|---|---|---|
| **CLI** | Git-compatible command surface with formatted output, JSON mode, and passthrough | Cobra |
| **TUI** | Lazygit-style multi-panel layout, interactive staging, commit graph, and themes | Bubble Tea |
| **Web UI** | Embedded offline-first web interface and PWA served on loopback | React + Vite |
| **Local AI** | Commit generation, diff explanation, and code review with rule-based fallback | Pluggable Engine |
| **Plugin System** | Sandboxed capabilities-based plugin runtime | wazero (pure Go) |
| **Forge Integration** | PR/MR management, issues, and CI status checks | GitHub + GitLab |

---

## Installation

### From Source

Ensure Go 1.23+ is installed:

```bash
git clone https://github.com/phantom1785322-hub/gitforge.git
cd gitforge
make build
./bin/gitforge version
```

### Pre-built Binaries & Package Managers

Pre-compiled static binaries for all architectures and Homebrew / Scoop / AUR packages are coming soon with the v0.1.0 release.

---

## Quick Start

```bash
# Check repository status
gitforge status

# Stage changes
gitforge add .

# Create a commit
gitforge commit -m "feat: initial commit"

# View graphical commit history
gitforge log --graph --oneline

# Launch the interactive terminal UI
gitforge tui

# Launch the local web UI
gitforge web

# Run diagnostics
gitforge doctor
```

---

## Key Commands

| Command | Description |
|---|---|
| `gitforge status` | Show working tree status (`--json` supported) |
| `gitforge log` | View commit history (`--oneline`, `--graph`, `--json`) |
| `gitforge diff` | Inspect changes between commits, index, or working tree |
| `gitforge show [ref]` | Show detailed commit metadata, diff, and diffstat |
| `gitforge add [paths]` | Stage files (supports `-p` hunk-level staging) |
| `gitforge commit` | Record changes to the repository (`-m`, `--amend`) |
| `gitforge branch` | List, create, or delete branches |
| `gitforge checkout` | Switch branches or restore working tree files |
| `gitforge remote` | Manage tracked repositories |
| `gitforge fetch` / `pull` / `push` | Synchronize commits with remotes |
| `gitforge stash` | Stash changes in a dirty working directory |
| `gitforge tag` | Create, list, or delete tags |
| `gitforge config` | Inspect or update configuration (secrets automatically redacted) |
| `gitforge ai commit` | Generate commit messages using local AI |
| `gitforge pr list` | View and review pull requests from GitHub / GitLab |
| `gitforge plugin list` | Manage sandboxed WASM plugins |
| `gitforge doctor` | Verify system health, git version, terminal, and AI setup |
| `gitforge <any>` | Transparently falls through to standard git |

---

## TUI Keybindings

The interactive TUI provides responsive keyboard-driven navigation:

| Key | Action |
|---|---|
| `Tab` | Switch focus between panels |
| `j` / `k` (or arrows) | Move down / up |
| `Enter` | Select / confirm / open detail |
| `Space` | Toggle stage / unstage file or hunk |
| `s` | Stage selected item |
| `c` | Open commit dialog |
| `p` | Push to current remote branch |
| `P` | Pull from remote |
| `r` | Refresh repository state |
| `f` | Fetch from remotes |
| `d` | Show diff for selected item |
| `l` | Show commit log / graph view |
| `b` | Open branch switcher |
| `?` | Toggle keymap help overlay |
| `q` | Quit GitForge |

---

## Architecture

GitForge follows strict inward-pointing layering:

```
┌──────────────────────────────────────────────┐
│           Front Ends (CLI / TUI / Web)       │
└──────────────────────┬───────────────────────┘
                       │
               ┌───────▼────────┐
               │  core.Engine   │  ← Facade & Event Bus
               └───────┬────────┘
     ┌─────────────────┼──────────────────┐
     │         │       │        │         │
┌────▼────┐ ┌──▼───┐ ┌─▼───┐ ┌──▼───┐ ┌───▼────┐
│Git layer│ │Config│ │Store│ │  AI  │ │Plugins │
└─────────┘ └──────┘ └─────┘ └──────┘ └────────┘
```

For complete architectural details, see [docs/architecture.md](docs/architecture.md).

---

## Configuration

GitForge loads configuration with strict precedence:

`Environment Variables` → `./gitforge.toml` → `~/.config/gitforge/config.toml` → `/etc/gitforge/config.toml` → `Built-in Defaults`

See [gitforge.toml.example](gitforge.toml.example) for a fully documented configuration template.

---

## Contributing & License

Contributions are welcome! Please review [CONTRIBUTING.md](CONTRIBUTING.md) and our [SECURITY.md](SECURITY.md) policy.

GitForge is open-source software licensed under the [MIT License](LICENSE).
