package config

import "fmt"

// BuiltinThemes lists standard theme names supported out of the box.
var BuiltinThemes = map[string]bool{
	"auto":             true,
	"catppuccin-mocha": true,
	"dracula":          true,
	"nord":             true,
	"tokyo-night":      true,
	"gruvbox":          true,
	"everforest":       true,
	"rose-pine":        true,
}

// ValidIcons lists allowed icon set identifiers.
var ValidIcons = map[string]bool{
	"emoji": true,
	"nerd":  true,
	"ascii": true,
}

// ValidEngines lists supported AI inference engine types.
var ValidEngines = map[string]bool{
	"local": true,
	"http":  true,
	"wasm":  true,
}

// Validate checks the configuration for semantic correctness.
// If any key contains invalid values, it returns a descriptive *ConfigError.
func Validate(cfg *Config) error {
	if cfg.UI.Icons != "" && !ValidIcons[cfg.UI.Icons] {
		return &ConfigError{
			Key:     "ui.icons",
			Message: fmt.Sprintf("invalid icon set %q (must be one of: emoji, nerd, ascii)", cfg.UI.Icons),
		}
	}

	if cfg.AI.Engine != "" && !ValidEngines[cfg.AI.Engine] {
		return &ConfigError{
			Key:     "ai.engine",
			Message: fmt.Sprintf("invalid AI engine %q (must be one of: local, http, wasm)", cfg.AI.Engine),
		}
	}

	if cfg.Core.AutoFetchInterval < 0 {
		return &ConfigError{
			Key:     "core.auto_fetch_interval",
			Message: "auto_fetch_interval must be non-negative",
		}
	}

	if cfg.AI.ContextLines < 0 {
		return &ConfigError{
			Key:     "ai.context_lines",
			Message: "context_lines must be non-negative",
		}
	}

	return nil
}
