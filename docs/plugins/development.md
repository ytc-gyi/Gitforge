# Plugin Development Guide

> **Note:** The plugin system is planned for Phase 5 of GitForge development. This document provides an initial overview and will be expanded significantly as the plugin architecture is implemented.

## Overview

GitForge's plugin system will allow you to extend GitForge with custom functionality. Plugins can add new commands, modify behavior via hooks, integrate with external services, and customize the UI.

## Plugin Types

### Lua Plugins

Lightweight scripts for simple extensions. Lua plugins run in a sandboxed environment with access to a GitForge API.

```lua
-- Example: A pre-commit hook that checks for TODO comments
gitforge.hook("pre-commit", function(ctx)
    local files = ctx.staged_files()
    for _, file in ipairs(files) do
        local content = ctx.read_file(file)
        if content:find("TODO") then
            ctx.warn("File " .. file .. " contains TODO comments")
        end
    end
end)
```

### WASM Plugins

Compiled plugins that run in a WebAssembly sandbox. WASM plugins can be written in any language that compiles to WASM (Rust, Go, C, etc.) and offer near-native performance.

### Go Plugins

Native Go plugins for maximum integration. Go plugins have full access to GitForge's internal APIs.

## Plugin Structure

```
my-plugin/
├── plugin.toml          # Plugin metadata and configuration
├── main.lua             # Entry point (for Lua plugins)
├── main.wasm            # Entry point (for WASM plugins)
└── README.md            # Plugin documentation
```

### plugin.toml

```toml
[plugin]
name = "my-plugin"
version = "0.1.0"
description = "A brief description of what this plugin does"
author = "Your Name"
license = "MIT"
type = "lua"  # "lua", "wasm", or "go"

[permissions]
filesystem = ["read"]   # "read", "write"
network = false          # Does this plugin need network access?
exec = false             # Does this plugin need to execute external commands?
```

## Plugin API

> **Coming soon.** The plugin API is under design. Key areas will include:

- **Repository access** — Read objects, refs, index, and working tree.
- **Hook registration** — Register callbacks for Git lifecycle events.
- **Command registration** — Add new subcommands to the CLI.
- **UI extensions** — Add panels or views to the TUI.
- **Configuration** — Read and write plugin-specific configuration.
- **Logging** — Structured logging with level control.

## Security

- Plugins run in sandboxed environments with declared permissions.
- Unsigned plugins are blocked by default (configurable via `plugins.allow_unsigned`).
- Plugin signing uses Ed25519 keys.
- Network and filesystem access must be explicitly declared in `plugin.toml`.

## Installation

```bash
# Install a plugin from a URL (coming soon)
gitforge plugin install https://example.com/my-plugin.tar.gz

# Install a plugin from a local directory (coming soon)
gitforge plugin install ./my-plugin/

# List installed plugins (coming soon)
gitforge plugin list

# Remove a plugin (coming soon)
gitforge plugin remove my-plugin
```

## Next Steps

This guide will be expanded with:

- Complete API reference.
- Tutorial: Building your first plugin.
- Plugin testing framework.
- Plugin distribution and registry.
- Advanced topics: async operations, UI integration, inter-plugin communication.
