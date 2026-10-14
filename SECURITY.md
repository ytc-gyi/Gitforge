# Security Policy

## Supported Versions

| Version | Supported |
|---|---|
| Latest release | ✅ |
| Previous minor | ✅ (security fixes only) |
| Older versions | ❌ |

## Reporting Vulnerabilities

If you discover a security vulnerability in GitForge, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please:

1. Email the maintainers with a detailed description of the vulnerability.
2. Include steps to reproduce the issue.
3. Include the affected version(s).
4. Allow reasonable time for a fix before public disclosure.

We will acknowledge your report within 48 hours and aim to release a fix within 7 days for critical vulnerabilities.

## Security Design Principles

GitForge is designed with security and privacy as core tenets:

### No Telemetry

- GitForge does **not** collect, transmit, or store any telemetry data.
- No usage analytics, crash reports, or tracking of any kind.
- No network requests are made unless explicitly initiated by the user (e.g., push, pull, fetch, forge API calls).

### Local-First

- All operations run locally on your machine by default.
- AI features use local models — no data is sent to external AI services unless you explicitly configure an external endpoint.
- Configuration and data are stored locally in standard file formats.

### Secrets Redaction

- GitForge redacts secrets (tokens, passwords, keys) from all log output and error messages.
- Credentials are never written to disk in plaintext outside of user-controlled configuration files.
- Environment variable references (e.g., `${GITHUB_TOKEN}`) are supported to avoid hardcoding secrets.

### Additional Principles

- **Minimal permissions** — GitForge requests only the permissions it needs.
- **No eval / dynamic code execution** — No `eval`-style execution of untrusted input.
- **Dependency minimalism** — We keep the dependency tree small and auditable.
- **CGO disabled** — Builds use `CGO_ENABLED=0` for static binaries with no C dependencies.

## Security Best Practices for Users

- Store tokens using environment variables, not directly in configuration files.
- Keep GitForge updated to the latest version.
- Review plugin code before installing third-party plugins.
- Use HTTPS for remote URLs when possible.
