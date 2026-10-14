package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/phantom1785322-hub/gitforge/internal/testutil"
)

// captureOutput redirects os.Stdout and runs fn.
func captureOutput(fn func() error) (string, error) {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outChan <- buf.String()
	}()

	fnErr := fn()

	_ = w.Close()
	os.Stdout = origStdout
	out := <-outChan
	_ = r.Close()

	return out, fnErr
}

func TestCLI_Commands(t *testing.T) {
	repoDir := testutil.InitRepo(t)
	testutil.CommitFile(t, repoDir, "file.txt", "line 1\n", "Initial commit")

	backend, err := git.NewCLIBackend(repoDir)
	if err != nil {
		t.Fatal(err)
	}
	SetBackend(backend)

	// 1. version
	out, err := captureOutput(func() error {
		return ExecuteArgs([]string{"version"})
	})
	if err != nil || !strings.Contains(out, "gitforge version") {
		t.Errorf("version output mismatch: %v, out: %q", err, out)
	}

	// 2. doctor
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"doctor"})
	})
	if err != nil || !strings.Contains(out, "GitForge Doctor") {
		t.Errorf("doctor output mismatch: %v, out: %q", err, out)
	}

	// 3. status clean
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"status"})
	})
	if err != nil || !strings.Contains(out, "working tree clean") {
		t.Errorf("status output mismatch: %v, out: %q", err, out)
	}

	// 4. status --json
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"status", "--json"})
	})
	if err != nil {
		t.Errorf("status --json failed: %v", err)
	}

	// 5. log
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"log"})
	})
	if err != nil || !strings.Contains(out, "Initial commit") {
		t.Errorf("log output mismatch: %v, out: %q", err, out)
	}

	// 6. log --oneline
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"log", "--oneline"})
	})
	if err != nil || !strings.Contains(out, "Initial commit") {
		t.Errorf("log --oneline mismatch: %v, out: %q", err, out)
	}

	// 7. log --graph
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"log", "--graph"})
	})
	if err != nil || !strings.Contains(out, "*") {
		t.Errorf("log --graph mismatch: %v, out: %q", err, out)
	}

	// 8. diff (make a modification)
	if err := os.WriteFile(filepath.Join(repoDir, "file.txt"), []byte("line 1\nline 2\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"diff"})
	})
	if err != nil || !strings.Contains(out, "+line 2") {
		t.Errorf("diff mismatch: %v, out: %q", err, out)
	}

	// 9. show
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"show", "HEAD"})
	})
	if err != nil || !strings.Contains(out, "Initial commit") {
		t.Errorf("show mismatch: %v, out: %q", err, out)
	}

	// 10. config
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"config"})
	})
	if err != nil || !strings.Contains(out, "GitForge Configuration") {
		t.Errorf("config mismatch: %v, out: %q", err, out)
	}

	// 11. alias expansion: "st" expands to "status"
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"st"})
	})
	if err != nil || !strings.Contains(out, "file.txt") {
		t.Errorf("alias 'st' failed: %v, out: %q", err, out)
	}

	// 12. add & commit
	if err := os.WriteFile(filepath.Join(repoDir, "extra.txt"), []byte("extra\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"add", "extra.txt"})
	})
	if err != nil {
		t.Errorf("add failed: %v", err)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"commit", "-m", "add extra"})
	})
	if err != nil || !strings.Contains(out, "add extra") {
		t.Errorf("commit failed: %v, out: %q", err, out)
	}

	// 13. branch create, list, and delete
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"branch", "feature-x"})
	})
	if err != nil || !strings.Contains(out, "Created branch") {
		t.Errorf("branch create failed: %v, out: %q", err, out)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"branch"})
	})
	if err != nil || !strings.Contains(out, "feature-x") {
		t.Errorf("branch list failed: %v, out: %q", err, out)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"branch", "-d", "feature-x"})
	})
	if err != nil || !strings.Contains(out, "Deleted branch") {
		t.Errorf("branch delete failed: %v, out: %q", err, out)
	}

	// 14. checkout -b and checkout
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"checkout", "-b", "dev"})
	})
	if err != nil || !strings.Contains(out, "Switched to a new branch") {
		t.Errorf("checkout -b failed: %v, out: %q", err, out)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"checkout", "main"})
	})
	if err != nil || !strings.Contains(out, "Switched to branch") {
		t.Errorf("checkout main failed: %v, out: %q", err, out)
	}

	// 15. tag create, list, and delete
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"tag", "v0.1.0"})
	})
	if err != nil || !strings.Contains(out, "Created tag") {
		t.Errorf("tag create failed: %v, out: %q", err, out)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"tag"})
	})
	if err != nil || !strings.Contains(out, "v0.1.0") {
		t.Errorf("tag list failed: %v, out: %q", err, out)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"tag", "-d", "v0.1.0"})
	})
	if err != nil || !strings.Contains(out, "Deleted tag") {
		t.Errorf("tag delete failed: %v, out: %q", err, out)
	}

	// 16. remote add and list
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"remote", "add", "upstream", "https://github.com/upstream/repo.git"})
	})
	if err != nil || !strings.Contains(out, "Added remote") {
		t.Errorf("remote add failed: %v, out: %q", err, out)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"remote"})
	})
	if err != nil || !strings.Contains(out, "upstream") {
		t.Errorf("remote list failed: %v, out: %q", err, out)
	}

	// 17. stash push, list, pop
	if err := os.WriteFile(filepath.Join(repoDir, "stash.txt"), []byte("stashing\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_ = ExecuteArgs([]string{"add", "stash.txt"})
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"stash"})
	})
	if err != nil || !strings.Contains(out, "Saved working directory") {
		t.Errorf("stash push failed: %v, out: %q", err, out)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"stash", "list"})
	})
	if err != nil || !strings.Contains(out, "stash@{0}") {
		t.Errorf("stash list failed: %v, out: %q", err, out)
	}
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"stash", "pop"})
	})
	if err != nil || !strings.Contains(out, "Applied and dropped stash") {
		t.Errorf("stash pop failed: %v, out: %q", err, out)
	}

	// 18. init a new repository in a separate directory
	initTarget := filepath.Join(t.TempDir(), "cli-init-repo")
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"init", initTarget})
	})
	if err != nil || !strings.Contains(out, "Initialized empty Git repository") {
		t.Errorf("init failed: %v, out: %q", err, out)
	}

	// 19. shell completions (bash, zsh, fish, powershell)
	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		out, err = captureOutput(func() error {
			return ExecuteArgs([]string{"completion", shell})
		})
		if err != nil || len(out) == 0 {
			t.Errorf("completion %s failed: %v, len: %d", shell, err, len(out))
		}
	}

	// 20. passthrough of unknown commands to git (e.g. rev-parse)
	out, err = captureOutput(func() error {
		return ExecuteArgs([]string{"rev-parse", "--is-inside-work-tree"})
	})
	if err != nil || !strings.Contains(out, "true") {
		t.Errorf("passthrough failed: %v, out: %q", err, out)
	}
}
