package git

import (
	"testing"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

func TestParseDiff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantFiles int
		checkDiff func(t *testing.T, res *core.DiffResult)
	}{
		{
			name:      "empty diff",
			input:     "",
			wantFiles: 0,
		},
		{
			name: "modified file with hunks",
			input: `diff --git a/main.go b/main.go
index 1234567..89abcdef 100644
--- a/main.go
+++ b/main.go
@@ -1,4 +1,5 @@
 package main
 
-func old() {}
+func newFunc() {}
+func added() {}
`,
			wantFiles: 1,
			checkDiff: func(t *testing.T, res *core.DiffResult) {
				t.Helper()
				f := res.Files[0]
				if f.OldName != "main.go" || f.NewName != "main.go" {
					t.Errorf("file name mismatch: %s -> %s", f.OldName, f.NewName)
				}
				if f.Status != core.DiffStatusModified {
					t.Errorf("expected modified, got %v", f.Status)
				}
				if len(f.Hunks) != 1 {
					t.Fatalf("expected 1 hunk, got %d", len(f.Hunks))
				}
				h := f.Hunks[0]
				if h.OldStart != 1 || h.OldLines != 4 || h.NewStart != 1 || h.NewLines != 5 {
					t.Errorf("hunk line ranges mismatch: %+v", h)
				}
				// Lines: context(package main), context(), deleted(-func old()), added(+func newFunc()), added(+func added())
				if len(h.Lines) < 4 {
					t.Fatalf("expected at least 4 lines in hunk, got %d", len(h.Lines))
				}
			},
		},
		{
			name: "new file",
			input: `diff --git a/hello.txt b/hello.txt
new file mode 100644
index 0000000..ce01362
--- /dev/null
+++ b/hello.txt
@@ -0,0 +1,2 @@
+hello
+world
`,
			wantFiles: 1,
			checkDiff: func(t *testing.T, res *core.DiffResult) {
				t.Helper()
				f := res.Files[0]
				if f.NewName != "hello.txt" {
					t.Errorf("file name mismatch: %s", f.NewName)
				}
				if f.Status != core.DiffStatusAdded {
					t.Errorf("expected DiffStatusAdded, got %v", f.Status)
				}
				if len(f.Hunks) != 1 || len(f.Hunks[0].Lines) != 2 {
					t.Errorf("unexpected hunks: %+v", f.Hunks)
				}
			},
		},
		{
			name: "deleted file",
			input: `diff --git a/bye.txt b/bye.txt
deleted file mode 100644
index ce01362..0000000
--- a/bye.txt
+++ /dev/null
@@ -1,1 +0,0 @@
-goodbye
`,
			wantFiles: 1,
			checkDiff: func(t *testing.T, res *core.DiffResult) {
				t.Helper()
				f := res.Files[0]
				if f.OldName != "bye.txt" {
					t.Errorf("file name mismatch: %s", f.OldName)
				}
				if f.Status != core.DiffStatusDeleted {
					t.Errorf("expected DiffStatusDeleted, got %v", f.Status)
				}
			},
		},
		{
			name: "binary file",
			input: `diff --git a/image.png b/image.png
new file mode 100644
index 0000000..abcdef1
Binary files /dev/null and b/image.png differ
`,
			wantFiles: 1,
			checkDiff: func(t *testing.T, res *core.DiffResult) {
				t.Helper()
				f := res.Files[0]
				if !f.Binary {
					t.Error("expected binary file to be marked as Binary")
				}
			},
		},
		{
			name: "renamed file",
			input: `diff --git a/old.go b/new.go
similarity index 100%
rename from old.go
rename to new.go
`,
			wantFiles: 1,
			checkDiff: func(t *testing.T, res *core.DiffResult) {
				t.Helper()
				f := res.Files[0]
				if f.OldName != "old.go" || f.NewName != "new.go" {
					t.Errorf("names mismatch: %s -> %s", f.OldName, f.NewName)
				}
				if f.Status != core.DiffStatusRenamed {
					t.Errorf("expected DiffStatusRenamed, got %v", f.Status)
				}
			},
		},
		{
			name: "copied file",
			input: `diff --git a/src.go b/dst.go
similarity index 90%
copy from src.go
copy to dst.go
`,
			wantFiles: 1,
			checkDiff: func(t *testing.T, res *core.DiffResult) {
				t.Helper()
				f := res.Files[0]
				if f.OldName != "src.go" || f.NewName != "dst.go" {
					t.Errorf("copied names mismatch: %s -> %s", f.OldName, f.NewName)
				}
				if f.Status != core.DiffStatusCopied {
					t.Errorf("expected DiffStatusCopied, got %v", f.Status)
				}
			},
		},
		{
			name: "mode change",
			input: `diff --git a/script.sh b/script.sh
old mode 100644
new mode 100755
`,
			wantFiles: 1,
			checkDiff: func(t *testing.T, res *core.DiffResult) {
				t.Helper()
				f := res.Files[0]
				if f.OldMode != "100644" || f.NewMode != "100755" {
					t.Errorf("modes mismatch: %s -> %s", f.OldMode, f.NewMode)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := ParseDiff(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res.Files) != tt.wantFiles {
				t.Fatalf("expected %d files, got %d", tt.wantFiles, len(res.Files))
			}
			if tt.checkDiff != nil {
				tt.checkDiff(t, res)
			}
		})
	}
}

func TestParseDiffStat(t *testing.T) {
	t.Parallel()

	input := ` main.go | 10 +++++-----
 doc.go  |  2 ++
 2 files changed, 7 insertions(+), 5 deletions(-)
`
	stats := ParseDiffStat(input)
	if len(stats) != 2 {
		t.Fatalf("expected 2 file stats, got %d", len(stats))
	}
	if stats[0].File != "main.go" || stats[0].Additions != 5 || stats[0].Deletions != 5 {
		t.Errorf("stat 0 mismatch: %+v", stats[0])
	}
	if stats[1].File != "doc.go" || stats[1].Additions != 2 || stats[1].Deletions != 0 {
		t.Errorf("stat 1 mismatch: %+v", stats[1])
	}
}
