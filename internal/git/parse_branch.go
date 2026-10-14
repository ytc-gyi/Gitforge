package git

import (
	"strconv"
	"strings"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

// ParseBranches parses `git branch --format` output with NUL-separated fields.
// Format: %(refname:short)%00%(objectname:short)%00%(upstream:short)%00%(upstream:track,nobracket)%00%(HEAD)
func ParseBranches(data string) ([]core.Branch, error) {
	if strings.TrimSpace(data) == "" {
		return nil, nil
	}

	var branches []core.Branch
	lines := strings.Split(strings.TrimSpace(data), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\x00")
		if len(fields) < 5 {
			continue
		}

		branch := core.Branch{
			Name:     fields[0],
			Hash:     fields[1],
			Upstream: fields[2],
			Current:  strings.TrimSpace(fields[4]) == "*",
		}

		// Parse ahead/behind from track info
		track := fields[3]
		if track != "" {
			branch.Ahead, branch.Behind = parseTrackInfo(track)
		}

		branches = append(branches, branch)
	}

	return branches, nil
}

func parseTrackInfo(track string) (ahead, behind int) {
	// Format: "ahead N", "behind N", "ahead N, behind M"
	parts := strings.Split(track, ", ")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "ahead ") {
			ahead, _ = strconv.Atoi(strings.TrimPrefix(part, "ahead "))
		} else if strings.HasPrefix(part, "behind ") {
			behind, _ = strconv.Atoi(strings.TrimPrefix(part, "behind "))
		}
	}
	return
}
