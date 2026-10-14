// Package config handles GitForge configuration loading, validation, and access.
// Configuration is loaded from multiple sources with a defined precedence order:
// environment variables -> ./gitforge.toml -> ~/.config/gitforge/config.toml (or OS equivalent)
// -> /etc/gitforge/config.toml -> built-in defaults.
package config
