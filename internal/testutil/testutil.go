package testutil

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// InitRepo creates a fresh temporary Git repository configured with user.name and user.email.
// It cleans up automatically when the test finishes.
func InitRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	MustGit(t, dir, "init", "-b", "main")
	MustGit(t, dir, "config", "user.name", "GitForge Test")
	MustGit(t, dir, "config", "user.email", "test@gitforge.dev")
	MustGit(t, dir, "config", "commit.gpgsign", "false")

	return dir
}

// CommitFile creates or overwrites a file, adds it to the index, and creates a commit.
// Returns the trimmed commit hash.
func CommitFile(t *testing.T, repoDir, relPath, content, msg string) string {
	t.Helper()

	fullPath := filepath.Join(repoDir, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", relPath, err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", relPath, err)
	}

	MustGit(t, repoDir, "add", relPath)
	MustGit(t, repoDir, "commit", "-m", msg)

	out := MustGit(t, repoDir, "rev-parse", "HEAD")
	return strings.TrimSpace(out)
}

// CreateBranch creates a new branch pointing to current HEAD.
func CreateBranch(t *testing.T, repoDir, branchName string) {
	t.Helper()
	MustGit(t, repoDir, "branch", branchName)
}

// Checkout switches to the specified branch or ref.
func Checkout(t *testing.T, repoDir, ref string) {
	t.Helper()
	MustGit(t, repoDir, "checkout", ref)
}

// MustGit runs a git command in the given repository directory, failing the test if it exits non-zero.
func MustGit(t *testing.T, repoDir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=GitForge Test",
		"GIT_AUTHOR_EMAIL=test@gitforge.dev",
		"GIT_COMMITTER_NAME=GitForge Test",
		"GIT_COMMITTER_EMAIL=test@gitforge.dev",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("git %s failed: %v\nstderr: %s", strings.Join(args, " "), err, stderr.String())
	}

	return stdout.String()
}
