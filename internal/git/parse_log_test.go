package git

import (
	"testing"
	"time"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

func TestParseLog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		wantCount   int
		checkCommit func(t *testing.T, commits []core.Commit)
	}{
		{
			name:      "empty log",
			input:     "",
			wantCount: 0,
		},
		{
			name: "single commit with refs",
			// Fields: %H\0%h\0%P\0%T\0%an\0%ae\0%aI\0%cn\0%ce\0%cI\0%s\0%b\0%D\x01
			input:     "1234567890abcdef1234567890abcdef12345678\x001234567\x00parent1 parent2\x00tree123\x00Alice\x00alice@example.com\x002026-10-14T10:00:00Z\x00Bob\x00bob@example.com\x002026-10-14T10:05:00Z\x00feat: add feature\x00Detailed body here.\x00HEAD -> main, origin/main, tag: v1.0.0\x01",
			wantCount: 1,
			checkCommit: func(t *testing.T, commits []core.Commit) {
				t.Helper()
				c := commits[0]
				if c.Hash != "1234567890abcdef1234567890abcdef12345678" {
					t.Errorf("hash mismatch: %s", c.Hash)
				}
				if c.AbbrevHash != "1234567" {
					t.Errorf("abbrev hash mismatch: %s", c.AbbrevHash)
				}
				if len(c.ParentHashes) != 2 || c.ParentHashes[0] != "parent1" || c.ParentHashes[1] != "parent2" {
					t.Errorf("parents mismatch: %+v", c.ParentHashes)
				}
				if c.Author.Name != "Alice" || c.Author.Email != "alice@example.com" {
					t.Errorf("author mismatch: %+v", c.Author)
				}
				if c.Author.Date.UTC() != time.Date(2026, 10, 14, 10, 0, 0, 0, time.UTC) {
					t.Errorf("author date mismatch: %v", c.Author.Date)
				}
				if c.Subject != "feat: add feature" {
					t.Errorf("subject mismatch: %s", c.Subject)
				}
				if c.Body != "Detailed body here." {
					t.Errorf("body mismatch: %s", c.Body)
				}
				if len(c.Refs) != 3 {
					t.Fatalf("expected 3 refs, got %d", len(c.Refs))
				}
				if c.Refs[0].Type != core.RefTypeHEAD || c.Refs[0].Name != "main" {
					t.Errorf("ref 0 mismatch: %+v", c.Refs[0])
				}
				if c.Refs[1].Type != core.RefTypeRemoteBranch || c.Refs[1].Name != "origin/main" || c.Refs[1].Remote != "origin" {
					t.Errorf("ref 1 mismatch: %+v", c.Refs[1])
				}
				if c.Refs[2].Type != core.RefTypeTag || c.Refs[2].Name != "v1.0.0" {
					t.Errorf("ref 2 mismatch: %+v", c.Refs[2])
				}
			},
		},
		{
			name: "multiple commits",
			input: "aaa\x00a\x00\x00t1\x00Author1\x00a1@test.com\x002026-01-01T00:00:00Z\x00Author1\x00a1@test.com\x002026-01-01T00:00:00Z\x00Commit 1\x00\x00\x01" +
				"bbb\x00b\x00aaa\x00t2\x00Author2\x00a2@test.com\x002026-01-02T00:00:00Z\x00Author2\x00a2@test.com\x002026-01-02T00:00:00Z\x00Commit 2\x00body text\x00\x01",
			wantCount: 2,
			checkCommit: func(t *testing.T, commits []core.Commit) {
				t.Helper()
				if commits[0].Hash != "aaa" || commits[1].Hash != "bbb" {
					t.Errorf("hashes mismatch: %s, %s", commits[0].Hash, commits[1].Hash)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			commits, err := ParseLog(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(commits) != tt.wantCount {
				t.Fatalf("expected %d commits, got %d", tt.wantCount, len(commits))
			}
			if tt.checkCommit != nil {
				tt.checkCommit(t, commits)
			}
		})
	}
}
