package cli

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

// IsTTY returns true if stdout is connected to a terminal and NO_COLOR is not set.
func IsTTY() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// Styling definitions with graceful plain-text fallback when not in a TTY.
var (
	StyleBold    = lipgloss.NewStyle().Bold(true)
	StyleMuted   = lipgloss.NewStyle().Faint(true)
	StyleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	StyleWarning = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	StyleError   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red
	StyleHash    = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	StyleBranch  = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // cyan
	StyleAdded   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	StyleDeleted = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red
)

// FormatHash formats a commit hash if TTY is available.
func FormatHash(hash string) string {
	if !IsTTY() {
		return hash
	}
	return StyleHash.Render(hash)
}

// FormatBranch formats a branch name if TTY is available.
func FormatBranch(name string) string {
	if !IsTTY() {
		return name
	}
	return StyleBranch.Render(name)
}

// FormatSuccess formats a success message.
func FormatSuccess(msg string) string {
	if !IsTTY() {
		return msg
	}
	return StyleSuccess.Render(msg)
}

// FormatError formats an error message.
func FormatError(msg string) string {
	if !IsTTY() {
		return msg
	}
	return StyleError.Render(msg)
}
