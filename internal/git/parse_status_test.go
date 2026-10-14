package git

import (
	"testing"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

func TestParseStatusPorcelainV2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		wantCount  int
		checkEntry func(t *testing.T, entries []core.StatusEntry)
		wantErr    bool
	}{
		{
			name:      "empty status",
			input:     "",
			wantCount: 0,
		},
		{
			name: "staged new file, modified worktree, untracked, ignored",
			// 1 <XY> <sub> <mH> <mI> <mW> <hH> <hI> <path>\0
			// ? <path>\0
			// ! <path>\0
			input: "1 A. N... 000000 100644 100644 0000000 1234567 staged.txt\x00" +
				"1 .M N... 100644 100644 100644 1234567 1234567 modified.txt\x00" +
				"? untracked.txt\x00" +
				"! ignored.log\x00",
			wantCount: 4,
			checkEntry: func(t *testing.T, entries []core.StatusEntry) {
				t.Helper()
				// Staged file
				if entries[0].Path != "staged.txt" || entries[0].Index != core.StatusAdded || entries[0].Worktree != core.StatusCode('.') {
					t.Errorf("entry 0 mismatch: %+v", entries[0])
				}
				if !entries[0].IsStaged() {
					t.Error("expected entry 0 to be staged")
				}

				// Modified worktree file
				if entries[1].Path != "modified.txt" || entries[1].Worktree != core.StatusModified {
					t.Errorf("entry 1 mismatch: %+v", entries[1])
				}
				if !entries[1].IsUnstaged() {
					t.Error("expected entry 1 to be unstaged")
				}

				// Untracked file
				if entries[2].Path != "untracked.txt" || entries[2].Index != core.StatusUntracked {
					t.Errorf("entry 2 mismatch: %+v", entries[2])
				}

				// Ignored file
				if entries[3].Path != "ignored.log" || entries[3].Index != core.StatusIgnored {
					t.Errorf("entry 3 mismatch: %+v", entries[3])
				}
			},
		},
		{
			name: "renamed file with -z",
			// 2 <XY> <sub> <mH> <mI> <mW> <hH> <hI> <X><score> <path>\0<origPath>\0
			input:     "2 R. N... 100644 100644 100644 1234567 1234567 R100 new_name.txt\x00old_name.txt\x00",
			wantCount: 1,
			checkEntry: func(t *testing.T, entries []core.StatusEntry) {
				t.Helper()
				if entries[0].Path != "new_name.txt" || entries[0].OrigPath != "old_name.txt" {
					t.Errorf("renamed entry mismatch: %+v", entries[0])
				}
				if entries[0].Index != core.StatusRenamed {
					t.Errorf("expected StatusRenamed, got %s", entries[0].Index)
				}
			},
		},
		{
			name: "unmerged conflict file",
			// u <XY> <sub> <m1> <m2> <m3> <mW> <h1> <h2> <h3> <path>\0
			input:     "u UU N... 100644 100644 100644 100644 1111111 2222222 3333333 conflict.txt\x00",
			wantCount: 1,
			checkEntry: func(t *testing.T, entries []core.StatusEntry) {
				t.Helper()
				if entries[0].Path != "conflict.txt" || entries[0].Index != core.StatusUnmerged {
					t.Errorf("unmerged entry mismatch: %+v", entries[0])
				}
			},
		},
		{
			name:      "headers only are skipped",
			input:     "# branch.oid abc1234\x00# branch.head main\x00",
			wantCount: 0,
		},
		{
			name:    "malformed ordinary entry",
			input:   "1 M\x00",
			wantErr: true,
		},
		{
			name:    "malformed renamed entry",
			input:   "2 R\x00old\x00",
			wantErr: true,
		},
		{
			name:    "malformed unmerged entry",
			input:   "u UU\x00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			entries, err := ParseStatusPorcelainV2(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseStatusPorcelainV2 error = %v, wantErr = %v", err, tt.wantErr)
			}
			if len(entries) != tt.wantCount {
				t.Fatalf("got %d entries, want %d", len(entries), tt.wantCount)
			}
			if tt.checkEntry != nil {
				tt.checkEntry(t, entries)
			}
		})
	}
}
