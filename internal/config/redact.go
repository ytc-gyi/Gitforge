package config

import "strings"

// RedactedPlaceholder is the string substituted for secret tokens.
const RedactedPlaceholder = "********"

var secretKeys = map[string]bool{
	"github.token":          true,
	"gitlab.token":          true,
	"gitforge_github_token": true,
	"gitforge_gitlab_token": true,
	"token":                 true,
}

// IsSecret returns true if the key name indicates a sensitive secret value.
func IsSecret(key string) bool {
	return secretKeys[strings.ToLower(key)]
}

// RedactedConfig returns a deep copy of cfg with sensitive values masked.
func RedactedConfig(cfg *Config) *Config {
	if cfg == nil {
		return nil
	}
	copyCfg := *cfg
	if copyCfg.GitHub.Token != "" {
		copyCfg.GitHub.Token = RedactedPlaceholder
	}
	if copyCfg.GitLab.Token != "" {
		copyCfg.GitLab.Token = RedactedPlaceholder
	}
	return &copyCfg
}

// RedactValue masks the value if the key corresponds to a secret.
func RedactValue(key, value string) string {
	if IsSecret(key) && value != "" {
		return RedactedPlaceholder
	}
	return value
}
