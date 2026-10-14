package core

import (
	"errors"
	"fmt"
)

var (
	// ErrNotARepo is returned when the current directory is not a git repository.
	ErrNotARepo = errors.New("not a git repository")

	// ErrGitNotFound is returned when the git binary cannot be found.
	ErrGitNotFound = errors.New("git executable not found in PATH")

	// ErrGitTooOld is returned when the git version is too old.
	ErrGitTooOld = errors.New("git version is too old (minimum: 2.20.0)")

	// ErrConfigNotFound is returned when a config file is not found.
	ErrConfigNotFound = errors.New("config file not found")

	// ErrConfigInvalid is returned when a config file has invalid content.
	ErrConfigInvalid = errors.New("invalid configuration")

	// ErrMergeConflict is returned when a merge conflict is detected.
	ErrMergeConflict = errors.New("merge conflict")

	// ErrDetachedHead is returned when HEAD is detached.
	ErrDetachedHead = errors.New("HEAD is detached")

	// ErrRefNotFound is returned when a ref cannot be found.
	ErrRefNotFound = errors.New("reference not found")

	// ErrPluginNotFound is returned when a plugin cannot be found.
	ErrPluginNotFound = errors.New("plugin not found")

	// ErrPluginDenied is returned when a plugin capability is denied.
	ErrPluginDenied = errors.New("plugin capability denied")

	// ErrModelNotFound is returned when an AI model is not found.
	ErrModelNotFound = errors.New("AI model not found")

	// ErrNoToken is returned when auth token is missing.
	ErrNoToken = errors.New("authentication token not configured")
)

// GitError wraps a git command failure with the command and stderr.
type GitError struct {
	Command string
	Args    []string
	Stderr  string
	Err     error
}

func (e *GitError) Error() string {
	return fmt.Sprintf("git %s failed: %s", e.Command, e.Stderr)
}

func (e *GitError) Unwrap() error {
	return e.Err
}

// ConfigError represents a configuration error with file and key context.
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
