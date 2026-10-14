package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	cfg := DefaultConfig()

	if cfg.UI.Theme != "auto" {
		t.Errorf("expected default theme 'auto', got %q", cfg.UI.Theme)
	}
	if cfg.UI.Icons != "ascii" {
		t.Errorf("expected default icons 'ascii', got %q", cfg.UI.Icons)
	}
	if cfg.AI.Engine != "local" {
		t.Errorf("expected default AI engine 'local', got %q", cfg.AI.Engine)
	}
	if !cfg.UI.Mouse {
		t.Error("expected mouse enabled by default")
	}
	if !cfg.UI.Unicode {
		t.Error("expected unicode enabled by default")
	}
	if cfg.Plugins.AllowUnsigned {
		t.Error("expected allow_unsigned to be false by default")
	}
	if cfg.Aliases["co"] != "checkout" {
		t.Errorf("expected alias 'co' = 'checkout', got %q", cfg.Aliases["co"])
	}
}

func TestLoadFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	content := `
[core]
editor = "vim"
auto_fetch = true
auto_fetch_interval = 600

[ui]
theme = "dracula"
icons = "nerd"

[ai]
enabled = false
engine = "http"
endpoint = "http://127.0.0.1:8080"
context_lines = 5

[github]
token = "ghp_test123"

[aliases]
st = "status --short"
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := DefaultConfig()
	if err := loadFile(cfgPath, cfg); err != nil {
		t.Fatalf("loadFile failed: %v", err)
	}

	if cfg.Core.Editor != "vim" {
		t.Errorf("expected editor 'vim', got %q", cfg.Core.Editor)
	}
	if !cfg.Core.AutoFetch {
		t.Error("expected auto_fetch to be true")
	}
	if cfg.Core.AutoFetchInterval != 600 {
		t.Errorf("expected interval 600, got %d", cfg.Core.AutoFetchInterval)
	}
	if cfg.UI.Theme != "dracula" {
		t.Errorf("expected theme 'dracula', got %q", cfg.UI.Theme)
	}
	if cfg.UI.Icons != "nerd" {
		t.Errorf("expected icons 'nerd', got %q", cfg.UI.Icons)
	}
	if cfg.AI.Engine != "http" {
		t.Errorf("expected AI engine 'http', got %q", cfg.AI.Engine)
	}
	if cfg.AI.ContextLines != 5 {
		t.Errorf("expected context_lines 5, got %d", cfg.AI.ContextLines)
	}
	if cfg.GitHub.Token != "ghp_test123" {
		t.Errorf("expected token 'ghp_test123', got %q", cfg.GitHub.Token)
	}
	if cfg.Aliases["st"] != "status --short" {
		t.Errorf("expected alias override, got %q", cfg.Aliases["st"])
	}
}

func TestEnvExpansion(t *testing.T) {
	t.Setenv("TEST_GH_TOKEN_EXPAND", "secret_expanded_token")
	t.Setenv("TEST_EDITOR", "nano")

	cfg := DefaultConfig()
	cfg.GitHub.Token = "${TEST_GH_TOKEN_EXPAND}"
	cfg.Core.Editor = "$TEST_EDITOR"
	ExpandEnvVars(cfg)

	if cfg.GitHub.Token != "secret_expanded_token" {
		t.Errorf("expected expanded token, got %q", cfg.GitHub.Token)
	}
	if cfg.Core.Editor != "nano" {
		t.Errorf("expected expanded editor, got %q", cfg.Core.Editor)
	}
}

func TestApplyEnvOverrides(t *testing.T) {
	t.Setenv("GITFORGE_THEME", "nord")
	t.Setenv("GITFORGE_ICONS", "nerd")
	t.Setenv("GITFORGE_AI_ENGINE", "http")
	t.Setenv("GITFORGE_AI_ENDPOINT", "http://127.0.0.1:9000")
	t.Setenv("GITFORGE_AI_MODEL", "codellama")
	t.Setenv("GITFORGE_GITHUB_TOKEN", "token_from_gitforge_env")
	t.Setenv("GITFORGE_GITLAB_TOKEN", "token_from_gitlab_env")

	cfg := DefaultConfig()
	ApplyEnvOverrides(cfg)

	if cfg.UI.Theme != "nord" {
		t.Errorf("expected theme 'nord' from env, got %q", cfg.UI.Theme)
	}
	if cfg.UI.Icons != "nerd" {
		t.Errorf("expected icons 'nerd' from env, got %q", cfg.UI.Icons)
	}
	if cfg.AI.Engine != "http" {
		t.Errorf("expected engine 'http' from env, got %q", cfg.AI.Engine)
	}
	if cfg.AI.Endpoint != "http://127.0.0.1:9000" {
		t.Errorf("expected endpoint from env, got %q", cfg.AI.Endpoint)
	}
	if cfg.AI.Model != "codellama" {
		t.Errorf("expected model from env, got %q", cfg.AI.Model)
	}
	if cfg.GitHub.Token != "token_from_gitforge_env" {
		t.Errorf("expected token from env, got %q", cfg.GitHub.Token)
	}
	if cfg.GitLab.Token != "token_from_gitlab_env" {
		t.Errorf("expected token from env, got %q", cfg.GitLab.Token)
	}
}

func TestValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		modify  func(*Config)
		wantErr bool
	}{
		{
			name:    "valid default config",
			modify:  func(c *Config) {},
			wantErr: false,
		},
		{
			name:    "invalid icons",
			modify:  func(c *Config) { c.UI.Icons = "unknown_icons" },
			wantErr: true,
		},
		{
			name:    "invalid AI engine",
			modify:  func(c *Config) { c.AI.Engine = "deepseek" },
			wantErr: true,
		},
		{
			name:    "negative auto_fetch_interval",
			modify:  func(c *Config) { c.Core.AutoFetchInterval = -1 },
			wantErr: true,
		},
		{
			name:    "negative context_lines",
			modify:  func(c *Config) { c.AI.ContextLines = -1 },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := DefaultConfig()
			tt.modify(cfg)
			err := Validate(cfg)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRedaction(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()
	cfg.GitHub.Token = "ghp_super_secret_token"
	cfg.GitLab.Token = "glpat-another_secret_token"

	redacted := RedactedConfig(cfg)
	if redacted.GitHub.Token != RedactedPlaceholder {
		t.Errorf("expected redacted GitHub token %q, got %q", RedactedPlaceholder, redacted.GitHub.Token)
	}
	if redacted.GitLab.Token != RedactedPlaceholder {
		t.Errorf("expected redacted GitLab token %q, got %q", RedactedPlaceholder, redacted.GitLab.Token)
	}
	if cfg.GitHub.Token != "ghp_super_secret_token" {
		t.Error("original config modified by RedactedConfig")
	}

	// nil check
	if RedactedConfig(nil) != nil {
		t.Error("expected nil for nil input")
	}

	// RedactValue
	if RedactValue("github.token", "secret123") != RedactedPlaceholder {
		t.Errorf("expected %q, got %q", RedactedPlaceholder, RedactValue("github.token", "secret123"))
	}
	if RedactValue("ui.theme", "nord") != "nord" {
		t.Errorf("expected 'nord', got %q", RedactValue("ui.theme", "nord"))
	}
}

func TestConfigError(t *testing.T) {
	t.Parallel()

	err1 := &ConfigError{File: "conf.toml", Key: "ui.theme", Message: "invalid"}
	if err1.Error() != `config error in conf.toml, key "ui.theme": invalid` {
		t.Errorf("unexpected error string: %s", err1.Error())
	}

	err2 := &ConfigError{File: "conf.toml", Message: "syntax error"}
	if err2.Error() != "config error in conf.toml: syntax error" {
		t.Errorf("unexpected error string: %s", err2.Error())
	}

	err3 := &ConfigError{Message: "generic error"}
	if err3.Error() != "config error: generic error" {
		t.Errorf("unexpected error string: %s", err3.Error())
	}
}

func TestLoadPrecedence(t *testing.T) {
	dir := t.TempDir()
	userCfg := filepath.Join(dir, "user.toml")
	localCfg := filepath.Join(dir, "local.toml")

	if err := os.WriteFile(userCfg, []byte("[ui]\ntheme = \"dracula\"\nicons = \"nerd\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(localCfg, []byte("[ui]\ntheme = \"nord\"\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := DefaultConfig()
	if err := loadFile(userCfg, cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.UI.Theme != "dracula" || cfg.UI.Icons != "nerd" {
		t.Errorf("user config mismatch: %+v", cfg.UI)
	}

	if err := loadFile(localCfg, cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.UI.Theme != "nord" {
		t.Errorf("expected local config to override theme to 'nord', got %q", cfg.UI.Theme)
	}
	// icons should remain nerd because local.toml didn't override it
	if cfg.UI.Icons != "nerd" {
		t.Errorf("expected icons to remain 'nerd', got %q", cfg.UI.Icons)
	}
}

func TestLoadInvalidToml(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	badFile := filepath.Join(dir, "bad.toml")
	if err := os.WriteFile(badFile, []byte("[[[invalid toml"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := DefaultConfig()
	err := loadFile(badFile, cfg)
	if err == nil {
		t.Error("expected error loading invalid TOML, got nil")
	}
	cfgErr, ok := err.(*ConfigError)
	if !ok {
		t.Errorf("expected *ConfigError, got %T", err)
	}
	if cfgErr.Unwrap() == nil {
		t.Error("expected wrapped error inside ConfigError")
	}
}

func TestConfigPaths(t *testing.T) {
	uPath := UserConfigPath()
	if uPath == "" {
		t.Error("expected non-empty UserConfigPath")
	}
	sPath := SystemConfigPath()
	if sPath == "" {
		t.Error("expected non-empty SystemConfigPath")
	}
}

func TestGetAndLoadedPaths(t *testing.T) {
	cfg := Get()
	if cfg == nil {
		t.Fatal("expected non-nil config from Get()")
	}
	_ = LoadedPaths()
}

func TestLoadGlobal(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config from Load()")
	}
	got := Get()
	if got == nil {
		t.Fatal("expected non-nil Get() after Load()")
	}
}
