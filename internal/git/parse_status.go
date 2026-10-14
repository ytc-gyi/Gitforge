package git

import (
	"fmt"
	"strings"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

// ParseStatusPorcelainV2 parses the output of `git status --porcelain=v2 -z`.
// The -z flag uses NUL bytes as field terminators.
func ParseStatusPorcelainV2(data string) ([]core.StatusEntry, error) {
	if data == "" {
		return nil, nil
	}

	var entries []core.StatusEntry
	fields := strings.Split(data, "\x00")

	i := 0
	for i < len(fields) {
		field := fields[i]
		if field == "" {
			i++
			continue
		}

		switch {
		case strings.HasPrefix(field, "1 "):
			// Ordinary changed entry: 1 <XY> <sub> <mH> <mI> <mW> <hH> <hI> <path>
			entry, err := parseOrdinaryEntry(field)
			if err != nil {
				return nil, err
			}
			entries = append(entries, entry)
			i++

		case strings.HasPrefix(field, "2 "):
			// Renamed/copied entry: 2 <XY> <sub> <mH> <mI> <mW> <hH> <hI> <X><score> <path>
			// With -z, the path and origPath are separate NUL-terminated fields
			var origPath string
			if i+1 < len(fields) {
				origPath = fields[i+1]
				i += 2
			} else {
				i++
			}
			entry, err := parseRenamedEntry(field, origPath)
			if err != nil {
				return nil, err
			}
			entries = append(entries, entry)

		case strings.HasPrefix(field, "u "):
			// Unmerged entry: u <XY> <sub> <m1> <m2> <m3> <mW> <h1> <h2> <h3> <path>
			entry, err := parseUnmergedEntry(field)
			if err != nil {
				return nil, err
			}
			entries = append(entries, entry)
			i++

		case strings.HasPrefix(field, "? "):
			// Untracked: ? <path>
			path := strings.TrimPrefix(field, "? ")
			entries = append(entries, core.StatusEntry{
				Index:    core.StatusUntracked,
				Worktree: core.StatusUntracked,
				Path:     path,
			})
			i++

		case strings.HasPrefix(field, "! "):
			// Ignored: ! <path>
			path := strings.TrimPrefix(field, "! ")
			entries = append(entries, core.StatusEntry{
				Index:    core.StatusIgnored,
				Worktree: core.StatusIgnored,
				Path:     path,
			})
			i++

		case strings.HasPrefix(field, "# "):
			// Header: skip
			i++

		default:
			i++
		}
	}

	return entries, nil
}

func parseOrdinaryEntry(line string) (core.StatusEntry, error) {
	// Format: 1 <XY> <sub> <mH> <mI> <mW> <hH> <hI> <path>
	parts := strings.SplitN(line, " ", 9)
	if len(parts) < 9 {
		return core.StatusEntry{}, fmt.Errorf("invalid ordinary status entry: %q", line)
	}

	xy := parts[1]
	if len(xy) < 2 {
		return core.StatusEntry{}, fmt.Errorf("invalid XY status: %q", xy)
	}

	return core.StatusEntry{
		Index:     core.StatusCode(xy[0]),
		Worktree:  core.StatusCode(xy[1]),
		Path:      parts[8],
		Submodule: parts[2],
	}, nil
}

func parseRenamedEntry(line string, origPath string) (core.StatusEntry, error) {
	// Format: 2 <XY> <sub> <mH> <mI> <mW> <hH> <hI> <X><score> <path>
	parts := strings.SplitN(line, " ", 10)
	if len(parts) < 10 {
		return core.StatusEntry{}, fmt.Errorf("invalid renamed status entry: %q", line)
	}

	xy := parts[1]
	if len(xy) < 2 {
		return core.StatusEntry{}, fmt.Errorf("invalid XY status: %q", xy)
	}

	return core.StatusEntry{
		Index:     core.StatusCode(xy[0]),
		Worktree:  core.StatusCode(xy[1]),
		Path:      parts[9],
		OrigPath:  origPath,
		Submodule: parts[2],
	}, nil
}

func parseUnmergedEntry(line string) (core.StatusEntry, error) {
	// Format: u <XY> <sub> <m1> <m2> <m3> <mW> <h1> <h2> <h3> <path>
	parts := strings.SplitN(line, " ", 11)
	if len(parts) < 11 {
		return core.StatusEntry{}, fmt.Errorf("invalid unmerged status entry: %q", line)
	}

	xy := parts[1]
	if len(xy) < 2 {
		return core.StatusEntry{}, fmt.Errorf("invalid XY status: %q", xy)
	}

	return core.StatusEntry{
		Index:     core.StatusCode(xy[0]),
		Worktree:  core.StatusCode(xy[1]),
		Path:      parts[10],
		Submodule: parts[2],
	}, nil
}
