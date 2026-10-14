package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

// Config represents the complete GitForge configuration.
type Config struct {
	Core    CoreConfig        `toml:"core"`
	UI      UIConfig          `toml:"ui"`
	AI      AIConfig          `toml:"ai"`
	GitHub  GitHubConfig      `toml:"github"`
	GitLab  GitLabConfig      `toml:"gitlab"`
	Plugins PluginsConfig     `toml:"plugins"`
	Aliases map[string]string `toml:"aliases"`
	Keys    map[string]string `toml:"keys"`
}

// CoreConfig holds core application settings.
type CoreConfig struct {
	Editor            string `toml:"editor"`
	Pager             string `toml:"pager"`
	AutoFetch         bool   `toml:"auto_fetch"`
	AutoFetchInterval int    `toml:"auto_fetch_interval"` // in seconds
}

// UIConfig holds terminal and web UI settings.
type UIConfig struct {
	Theme       string `toml:"theme"`
	Icons       string `toml:"icons"` // emoji, nerd, ascii
	Mouse       bool   `toml:"mouse"`
	Unicode     bool   `toml:"unicode"`
	CompactMode bool   `toml:"compact_mode"`
}

// AIConfig holds local and remote AI assistant settings.
type AIConfig struct {
	Enabled            bool   `toml:"enabled"`
	Engine             string `toml:"engine"`   // local, http, wasm
	Endpoint           string `toml:"endpoint"` // e.g. http://127.0.0.1:11434
	Model              string `toml:"model"`
	AutoDownloadModels bool   `toml:"auto_download_models"`
	ContextLines       int    `toml:"context_lines"`
}

// GitHubConfig holds GitHub integration settings.
type GitHubConfig struct {
	Token string `toml:"token"`
	Host  string `toml:"host"`
}

// GitLabConfig holds GitLab integration settings.
type GitLabConfig struct {
	Token string `toml:"token"`
	Host  string `toml:"host"`
}

// PluginsConfig holds WASM plugin settings.
type PluginsConfig struct {
	Directory      string   `toml:"directory"`
	AllowUnsigned  bool     `toml:"allow_unsigned"`
	TrustedSources []string `toml:"trusted_sources"`
}

// DefaultConfig returns a Config populated with standard defaults.
func DefaultConfig() *Config {
	return &Config{
		Core: CoreConfig{
			Editor:            "",
			Pager:             "",
			AutoFetch:         false,
			AutoFetchInterval: 300,
		},
		UI: UIConfig{
			Theme:       "auto",
			Icons:       "ascii",
			Mouse:       true,
			Unicode:     true,
			CompactMode: false,
		},
		AI: AIConfig{
			Enabled:            true,
			Engine:             "local",
			Endpoint:           "http://127.0.0.1:11434",
			Model:              "",
			AutoDownloadModels: false,
			ContextLines:       3,
		},
		GitHub: GitHubConfig{
			Host: "github.com",
		},
		GitLab: GitLabConfig{
			Host: "gitlab.com",
		},
		Plugins: PluginsConfig{
			Directory:     "",
			AllowUnsigned: false,
		},
		Aliases: map[string]string{
			"co": "checkout",
			"br": "branch",
			"st": "status",
			"lg": "log --graph --oneline",
			"cm": "commit",
		},
		Keys: make(map[string]string),
	}
}

var (
	globalConfig *Config
	globalMu     sync.RWMutex
	loadedPaths  []string
)

// Load loads configuration from all sources with proper precedence:
// Environment variables > ./gitforge.toml > user config > system config > defaults.
func Load() (*Config, error) {
	cfg := DefaultConfig()
	var loaded []string

	// 1. System config
	if sysPath := SystemConfigPath(); sysPath != "" {
		if err := loadFile(sysPath, cfg); err == nil {
			loaded = append(loaded, sysPath)
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("loading system config: %w", err)
		}
	}

	// 2. User config (XDG / OS standard)
	if userPath := UserConfigPath(); userPath != "" {
		if err := loadFile(userPath, cfg); err == nil {
			loaded = append(loaded, userPath)
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("loading user config: %w", err)
		}
	}

	// 3. Local repository config (./gitforge.toml)
	localPath := "gitforge.toml"
	if err := loadFile(localPath, cfg); err == nil {
		loaded = append(loaded, localPath)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading local config: %w", err)
	}

	// 4. Environment variable overrides
	ApplyEnvOverrides(cfg)

	// 5. Expand env vars in string values (${VAR})
	ExpandEnvVars(cfg)

	// 6. Validate
	if err := Validate(cfg); err != nil {
		return nil, err
	}

	globalMu.Lock()
	globalConfig = cfg
	loadedPaths = loaded
	globalMu.Unlock()

	return cfg, nil
}

// Get returns the current active configuration. If Load has not been called,
// DefaultConfig is returned.
func Get() *Config {
	globalMu.RLock()
	defer globalMu.RUnlock()
	if globalConfig == nil {
		return DefaultConfig()
	}
	return globalConfig
}

// LoadedPaths returns paths of configuration files that were loaded.
func LoadedPaths() []string {
	globalMu.RLock()
	defer globalMu.RUnlock()
	copied := make([]string, len(loadedPaths))
	copy(copied, loadedPaths)
	return copied
}

func loadFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return &ConfigError{
			File:    path,
			Message: err.Error(),
			Err:     err,
		}
	}
	return nil
}

// UserConfigPath returns the platform-specific path to the user configuration file.
func UserConfigPath() string {
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "gitforge", "config.toml")
		}
		if home, _ := os.UserHomeDir(); home != "" {
			return filepath.Join(home, "AppData", "Roaming", "gitforge", "config.toml")
		}
	case "darwin":
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, "gitforge", "config.toml")
		}
		if home, _ := os.UserHomeDir(); home != "" {
			return filepath.Join(home, "Library", "Application Support", "gitforge", "config.toml")
		}
	default: // Linux, BSD, Termux
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, "gitforge", "config.toml")
		}
		if home, _ := os.UserHomeDir(); home != "" {
			return filepath.Join(home, ".config", "gitforge", "config.toml")
		}
	}
	return ""
}

// SystemConfigPath returns the platform-specific path to the system-wide configuration file.
func SystemConfigPath() string {
	switch runtime.GOOS {
	case "windows":
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = `C:\ProgramData`
		}
		return filepath.Join(programData, "gitforge", "config.toml")
	default:
		return "/etc/gitforge/config.toml"
	}
}

// ApplyEnvOverrides overrides configuration with environment variables if present.
func ApplyEnvOverrides(cfg *Config) {
	if v := os.Getenv("GITFORGE_THEME"); v != "" {
		cfg.UI.Theme = v
	}
	if v := os.Getenv("GITFORGE_ICONS"); v != "" {
		cfg.UI.Icons = v
	}
	if v := os.Getenv("GITFORGE_AI_ENGINE"); v != "" {
		cfg.AI.Engine = v
	}
	if v := os.Getenv("GITFORGE_AI_ENDPOINT"); v != "" {
		cfg.AI.Endpoint = v
	}
	if v := os.Getenv("GITFORGE_AI_MODEL"); v != "" {
		cfg.AI.Model = v
	}
	if v := os.Getenv("GITFORGE_GITHUB_TOKEN"); v != "" {
		cfg.GitHub.Token = v
	} else if v := os.Getenv("GITHUB_TOKEN"); v != "" && cfg.GitHub.Token == "" {
		cfg.GitHub.Token = v
	}
	if v := os.Getenv("GITFORGE_GITLAB_TOKEN"); v != "" {
		cfg.GitLab.Token = v
	} else if v := os.Getenv("GITLAB_TOKEN"); v != "" && cfg.GitLab.Token == "" {
		cfg.GitLab.Token = v
	}
}

// ExpandEnvVars expands ${VAR} or $VAR expressions in configuration string values.
func ExpandEnvVars(cfg *Config) {
	cfg.GitHub.Token = expandEnv(cfg.GitHub.Token)
	cfg.GitLab.Token = expandEnv(cfg.GitLab.Token)
	cfg.AI.Endpoint = expandEnv(cfg.AI.Endpoint)
	cfg.Plugins.Directory = expandEnv(cfg.Plugins.Directory)
	cfg.Core.Editor = expandEnv(cfg.Core.Editor)
	cfg.Core.Pager = expandEnv(cfg.Core.Pager)
}

func expandEnv(s string) string {
	if !strings.Contains(s, "$") {
		return s
	}
	return os.Expand(s, os.Getenv)
}

// ConfigError represents a configuration parsing or validation error.
type ConfigError struct {
	File    string
	Key     string
	Message string
	Err     error
}

func (e *ConfigError) Error() string {
	if e.File != "" && e.Key != "" {
		return fmt.Sprintf("config error in %s, key %q: %s", e.File, e.Key, e.Message)
	}
	if e.File != "" {
		return fmt.Sprintf("config error in %s: %s", e.File, e.Message)
	}
	return fmt.Sprintf("config error: %s", e.Message)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}
