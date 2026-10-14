# Contributing to GitForge

Thank you for your interest in contributing to GitForge! This guide will help you get started.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

---

## Development Setup

### Prerequisites

- **Go 1.23+** — [Install Go](https://go.dev/dl/)
- **golangci-lint** — Install via `make install-tools`
- **Git** — A working Git installation
- **Make** — GNU Make or compatible

### Getting Started

```bash
# Clone the repository
git clone https://github.com/phantom1785322-hub/gitforge.git
cd gitforge

# Install development tools
make install-tools

# Build the project
make build

# Run all checks (build, vet, lint, test)
make check
```

---

## Coding Standards

### Formatting

- All Go code must be formatted with `gofmt` and `goimports`.
- Run `make fmt` before committing.

### Linting

- All code must pass `golangci-lint` with the project configuration.
- Run `make lint` to check.

### Testing

- All new code must include tests.
- Run `make test` to execute the test suite.
- Tests must pass with the race detector enabled (`-race`).
- Aim for meaningful test coverage — focus on edge cases and error paths.

### General Guidelines

- Follow idiomatic Go style. When in doubt, consult [Effective Go](https://go.dev/doc/effective_go).
- Keep functions short and focused.
- Use meaningful variable and function names.
- Document exported types, functions, and packages.
- Handle all errors explicitly — never ignore returned errors.
- Avoid global mutable state.
- Use `context.Context` for cancellation and deadlines.

---

## Pull Request Process

1. **Fork** the repository and create a feature branch from `main`.
2. **Make your changes** in small, focused commits.
3. **Run `make check`** to ensure all checks pass locally.
4. **Push** your branch and open a Pull Request against `main`.
5. **Describe** your changes clearly in the PR description.
6. **Respond** to code review feedback promptly.

### PR Requirements

- All CI checks must pass.
- At least one maintainer approval is required.
- Squash-merge is preferred for single-feature PRs.

---

## Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description |
|---|---|
| `feat` | A new feature |
| `fix` | A bug fix |
| `docs` | Documentation changes |
| `style` | Code style changes (formatting, no logic change) |
| `refactor` | Code refactoring (no feature or fix) |
| `perf` | Performance improvements |
| `test` | Adding or updating tests |
| `build` | Build system or dependency changes |
| `ci` | CI configuration changes |
| `chore` | Maintenance tasks |

### Examples

```
feat(cli): add status command with color output
fix(core): handle symlinks in tree walker
docs(readme): update installation instructions
test(diff): add edge case tests for binary files
```

---

## Reporting Issues

- Use GitHub Issues to report bugs or request features.
- Include steps to reproduce for bugs.
- Include your OS, Go version, and GitForge version.

---

## Questions?

Open a GitHub Discussion or reach out in an issue. We're happy to help!
