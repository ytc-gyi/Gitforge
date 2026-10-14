package git

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

// ParseDiff parses unified diff output into structured DiffResult.
func ParseDiff(data string) (*core.DiffResult, error) {
	if strings.TrimSpace(data) == "" {
		return &core.DiffResult{}, nil
	}

	result := &core.DiffResult{}
	lines := strings.Split(data, "\n")

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Look for diff --git header
		if strings.HasPrefix(line, "diff --git ") {
			file, newI, err := parseDiffFile(lines, i)
			if err != nil {
				return nil, err
			}
			result.Files = append(result.Files, file)
			i = newI
		} else {
			i++
		}
	}

	return result, nil
}

func parseDiffFile(lines []string, start int) (core.DiffFile, int, error) {
	file := core.DiffFile{}
	i := start
	parsedStatus := false

	// Parse "diff --git a/path b/path"
	header := lines[i]
	parts := strings.SplitN(header, " ", 4)
	if len(parts) >= 4 {
		file.OldName = strings.TrimPrefix(parts[2], "a/")
		file.NewName = strings.TrimPrefix(parts[3], "b/")
	}
	i++

	// Parse extended headers (old mode, new mode, similarity, etc.)
	for i < len(lines) {
		line := lines[i]
		switch {
		case strings.HasPrefix(line, "old mode "):
			file.OldMode = strings.TrimPrefix(line, "old mode ")
			i++
		case strings.HasPrefix(line, "new mode "):
			file.NewMode = strings.TrimPrefix(line, "new mode ")
			i++
		case strings.HasPrefix(line, "new file mode "):
			file.Status = core.DiffStatusAdded
			file.NewMode = strings.TrimPrefix(line, "new file mode ")
			parsedStatus = true
			i++
		case strings.HasPrefix(line, "deleted file mode "):
			file.Status = core.DiffStatusDeleted
			file.OldMode = strings.TrimPrefix(line, "deleted file mode ")
			parsedStatus = true
			i++
		case strings.HasPrefix(line, "similarity index"):
			i++
		case strings.HasPrefix(line, "rename from "):
			file.OldName = strings.TrimPrefix(line, "rename from ")
			file.Status = core.DiffStatusRenamed
			parsedStatus = true
			i++
		case strings.HasPrefix(line, "rename to "):
			file.NewName = strings.TrimPrefix(line, "rename to ")
			i++
		case strings.HasPrefix(line, "copy from "):
			file.OldName = strings.TrimPrefix(line, "copy from ")
			file.Status = core.DiffStatusCopied
			parsedStatus = true
			i++
		case strings.HasPrefix(line, "copy to "):
			file.NewName = strings.TrimPrefix(line, "copy to ")
			i++
		case strings.HasPrefix(line, "index "):
			i++
		case strings.HasPrefix(line, "Binary files"):
			file.Binary = true
			i++
		case strings.HasPrefix(line, "--- "):
			// Start of unified diff, parse hunks
			goto parseHunks
		case strings.HasPrefix(line, "diff --git "):
			// Next file
			goto done
		default:
			i++
		}
	}
	goto done

parseHunks:
	// Skip --- line
	if i < len(lines) && strings.HasPrefix(lines[i], "--- ") {
		oldFile := strings.TrimPrefix(lines[i], "--- ")
		if oldFile == "/dev/null" && !parsedStatus {
			file.Status = core.DiffStatusAdded
			parsedStatus = true
		}
		i++
	}
	// Skip +++ line
	if i < len(lines) && strings.HasPrefix(lines[i], "+++ ") {
		newFile := strings.TrimPrefix(lines[i], "+++ ")
		if newFile == "/dev/null" && !parsedStatus {
			file.Status = core.DiffStatusDeleted
			parsedStatus = true
		}
		i++
	}

	// Parse hunks
	for i < len(lines) {
		line := lines[i]
		if strings.HasPrefix(line, "@@ ") {
			hunk, newI, err := parseDiffHunk(lines, i)
			if err != nil {
				return file, i, err
			}
			file.Hunks = append(file.Hunks, hunk)
			i = newI
		} else if strings.HasPrefix(line, "diff --git ") {
			break
		} else {
			i++
		}
	}

	if !parsedStatus && len(file.Hunks) > 0 {
		file.Status = core.DiffStatusModified
	}

done:
	return file, i, nil
}

func parseDiffHunk(lines []string, start int) (core.DiffHunk, int, error) {
	hunk := core.DiffHunk{}
	line := lines[start]

	// Parse @@ -oldStart,oldLines +newStart,newLines @@ header
	if err := parseHunkHeader(line, &hunk); err != nil {
		return hunk, start + 1, err
	}

	i := start + 1
	for i < len(lines) {
		line := lines[i]
		switch {
		case strings.HasPrefix(line, "+"):
			hunk.Lines = append(hunk.Lines, core.DiffLine{
				Type:    core.DiffLineAdded,
				Content: line[1:],
			})
			i++
		case strings.HasPrefix(line, "-"):
			hunk.Lines = append(hunk.Lines, core.DiffLine{
				Type:    core.DiffLineDeleted,
				Content: line[1:],
			})
			i++
		case strings.HasPrefix(line, " "):
			hunk.Lines = append(hunk.Lines, core.DiffLine{
				Type:    core.DiffLineContext,
				Content: line[1:],
			})
			i++
		case line == `\ No newline at end of file`:
			i++
		case strings.HasPrefix(line, "@@") || strings.HasPrefix(line, "diff --git"):
			// Next hunk or file
			return hunk, i, nil
		default:
			// Context line without leading space (empty line)
			if line == "" {
				// Could be end of diff or an empty context line
				// Check if next line is part of diff
				if i+1 < len(lines) {
					next := lines[i+1]
					if strings.HasPrefix(next, "+") || strings.HasPrefix(next, "-") ||
						strings.HasPrefix(next, " ") || strings.HasPrefix(next, "@@") {
						hunk.Lines = append(hunk.Lines, core.DiffLine{
							Type:    core.DiffLineContext,
							Content: "",
						})
						i++
						continue
					}
				}
				return hunk, i, nil
			}
			return hunk, i, nil
		}
	}

	return hunk, i, nil
}

func parseHunkHeader(line string, hunk *core.DiffHunk) error {
	// Format: @@ -oldStart,oldLines +newStart,newLines @@ optional header
	hunk.Header = line

	// Find the range info between @@ markers
	if !strings.HasPrefix(line, "@@ ") {
		return fmt.Errorf("invalid hunk header: %q", line)
	}

	endIdx := strings.Index(line[3:], " @@")
	if endIdx < 0 {
		return fmt.Errorf("invalid hunk header: %q", line)
	}
	rangeStr := line[3 : 3+endIdx]

	parts := strings.Split(rangeStr, " ")
	if len(parts) != 2 {
		return fmt.Errorf("invalid hunk range: %q", rangeStr)
	}

	// Parse old range (-start,lines)
	oldRange := strings.TrimPrefix(parts[0], "-")
	oldParts := strings.SplitN(oldRange, ",", 2)
	hunk.OldStart, _ = strconv.Atoi(oldParts[0])
	if len(oldParts) > 1 {
		hunk.OldLines, _ = strconv.Atoi(oldParts[1])
	} else {
		hunk.OldLines = 1
	}

	// Parse new range (+start,lines)
	newRange := strings.TrimPrefix(parts[1], "+")
	newParts := strings.SplitN(newRange, ",", 2)
	hunk.NewStart, _ = strconv.Atoi(newParts[0])
	if len(newParts) > 1 {
		hunk.NewLines, _ = strconv.Atoi(newParts[1])
	} else {
		hunk.NewLines = 1
	}

	return nil
}

// ParseDiffStat parses --stat output.
func ParseDiffStat(data string) []core.FileStat {
	var stats []core.FileStat
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "|") {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		fileName := strings.TrimSpace(parts[0])
		changeStr := strings.TrimSpace(parts[1])

		// Count + and - characters
		var adds, dels int
		for _, ch := range changeStr {
			switch ch {
			case '+':
				adds++
			case '-':
				dels++
			}
		}

		stats = append(stats, core.FileStat{
			File:      fileName,
			Additions: adds,
			Deletions: dels,
		})
	}
	return stats
}
