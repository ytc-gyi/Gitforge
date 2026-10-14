package git

import (
	"fmt"
	"strings"
	"time"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

// LogFormatParseable is the git log --format string used for machine-readable output.
// Field separator \x00 between fields within a commit, record separator \x01 between commits.
const LogFormatParseable = "%H%x00%h%x00%P%x00%T%x00%an%x00%ae%x00%aI%x00%cn%x00%ce%x00%cI%x00%s%x00%b%x00%D"

// ParseLog parses git log output in the custom format.
// Records are delimited by \x01 (record separator).
func ParseLog(data string) ([]core.Commit, error) {
	if strings.TrimSpace(data) == "" {
		return nil, nil
	}

	// Split on record separator
	records := strings.Split(data, "\x01")
	var commits []core.Commit

	for _, record := range records {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		commit, err := parseLogRecord(record)
		if err != nil {
			continue // skip malformed records
		}
		commits = append(commits, commit)
	}

	return commits, nil
}

func parseLogRecord(record string) (core.Commit, error) {
	fields := strings.Split(record, "\x00")
	// We expect 13 fields: Hash, AbbrevHash, Parents, Tree, AuthorName, AuthorEmail,
	// AuthorDate, CommitterName, CommitterEmail, CommitterDate, Subject, Body, Refs
	if len(fields) < 13 {
		return core.Commit{}, fmt.Errorf("expected at least 13 fields, got %d", len(fields))
	}

	authorDate, _ := time.Parse(time.RFC3339, strings.TrimSpace(fields[6]))
	committerDate, _ := time.Parse(time.RFC3339, strings.TrimSpace(fields[9]))

	var parents []string
	if parentStr := strings.TrimSpace(fields[2]); parentStr != "" {
		parents = strings.Split(parentStr, " ")
	}

	refs := parseRefs(strings.TrimSpace(fields[len(fields)-1]))

	return core.Commit{
		Hash:         strings.TrimSpace(fields[0]),
		AbbrevHash:   strings.TrimSpace(fields[1]),
		ParentHashes: parents,
		TreeHash:     strings.TrimSpace(fields[3]),
		Author: core.Person{
			Name:  strings.TrimSpace(fields[4]),
			Email: strings.TrimSpace(fields[5]),
			Date:  authorDate,
		},
		Committer: core.Person{
			Name:  strings.TrimSpace(fields[7]),
			Email: strings.TrimSpace(fields[8]),
			Date:  committerDate,
		},
		Subject: strings.TrimSpace(fields[10]),
		Body:    strings.TrimSpace(fields[11]),
		Refs:    refs,
	}, nil
}

func parseRefs(refStr string) []core.Ref {
	if refStr == "" {
		return nil
	}

	var refs []core.Ref
	parts := strings.Split(refStr, ", ")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		ref := core.Ref{Name: part}
		switch {
		case part == "HEAD":
			ref.Type = core.RefTypeHEAD
		case strings.HasPrefix(part, "HEAD -> "):
			ref.Type = core.RefTypeHEAD
			ref.Name = strings.TrimPrefix(part, "HEAD -> ")
		case strings.HasPrefix(part, "tag: "):
			ref.Type = core.RefTypeTag
			ref.Name = strings.TrimPrefix(part, "tag: ")
		case strings.Contains(part, "/"):
			ref.Type = core.RefTypeRemoteBranch
			slashIdx := strings.Index(part, "/")
			ref.Remote = part[:slashIdx]
		default:
			ref.Type = core.RefTypeBranch
		}
		refs = append(refs, ref)
	}

	return refs
}
