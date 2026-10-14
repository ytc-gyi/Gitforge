package core

import "fmt"

// VersionString returns the formatted version string.
func VersionString() string {
	return fmt.Sprintf("gitforge version %s (commit: %s, built: %s)", BuildVersion, BuildCommit, BuildDate)
}
