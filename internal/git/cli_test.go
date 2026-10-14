package git

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/phantom1785322-hub/gitforge/internal/core"
	"github.com/phantom1785322-hub/gitforge/internal/testutil"
)

func TestCLIBackend_Integration(t *testing.T) {
	repoDir := testutil.InitRepo(t)
	ctx := context.Background()

	backend, err := NewCLIBackend(repoDir)
	if err != nil {
		t.Fatalf("NewCLIBackend failed: %v", err)
	}

	// 1. RepoRoot
	root, err := backend.RepoRoot(ctx)
	if err != nil {
		t.Fatalf("RepoRoot failed: %v", err)
	}
	// On macOS/Linux symlinks might differ, compare canonical paths
	realRepo, _ := filepath.EvalSymlinks(repoDir)
	realRoot, _ := filepath.EvalSymlinks(root)
	if realRepo != realRoot {
		t.Errorf("RepoRoot mismatch: got %s, want %s", realRoot, realRepo)
	}

	// 2. Initial commit
	c1Hash := testutil.CommitFile(t, repoDir, "hello.txt", "hello world\n", "Initial commit")

	// 3. Version
	ver, err := backend.Version(ctx)
	if err != nil || ver == "" {
		t.Errorf("Version failed: %v, got %q", err, ver)
	}

	// 4. CurrentBranch
	branch, err := backend.CurrentBranch(ctx)
	if err != nil || branch != "main" {
		t.Errorf("CurrentBranch failed: %v, got %q", err, branch)
	}

	// 5. Log with options (MaxCount, Skip, Author, Grep, etc.)
	commits, err := backend.Log(ctx, LogOptions{MaxCount: 10, Skip: 0, Author: "GitForge", Grep: "Initial"})
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}
	if len(commits) != 1 || commits[0].Hash != c1Hash {
		t.Errorf("unexpected log result: %+v", commits)
	}

	// 5b. DiffRaw with various options
	rawDiff, err := backend.DiffRaw(ctx, DiffOptions{Stat: true, NameOnly: true, ContextLines: 3})
	if err != nil {
		t.Fatalf("DiffRaw failed: %v", err)
	}
	_ = rawDiff

	// 6. Show
	detail, err := backend.Show(ctx, "HEAD")
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}
	if detail.Hash != c1Hash || detail.Subject != "Initial commit" {
		t.Errorf("unexpected show detail: %+v", detail)
	}

	// 7. Status - clean
	entries, err := backend.Status(ctx, StatusOptions{})
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected clean status, got %d entries", len(entries))
	}

	// 8. Add untracked file and test status
	if err := os.WriteFile(filepath.Join(repoDir, "new.txt"), []byte("new content\n"), 0644); err != nil {
		t.Fatal(err)
	}
	entries, err = backend.Status(ctx, StatusOptions{})
	if err != nil {
		t.Fatalf("Status with untracked failed: %v", err)
	}
	if len(entries) != 1 || entries[0].Path != "new.txt" || entries[0].Index != core.StatusUntracked {
		t.Errorf("unexpected status for untracked: %+v", entries)
	}

	// 9. Add & Commit via backend
	if err := backend.Add(ctx, []string{"new.txt"}); err != nil {
		t.Fatalf("backend.Add failed: %v", err)
	}
	entries, err = backend.Status(ctx, StatusOptions{})
	if err != nil || len(entries) != 1 || !entries[0].IsStaged() {
		t.Errorf("expected staged status, got %+v", entries)
	}

	c2Hash, err := backend.Commit(ctx, CommitOptions{Message: "Second commit"})
	if err != nil {
		t.Fatalf("backend.Commit failed: %v", err)
	}
	if c2Hash == "" {
		t.Error("expected non-empty commit hash")
	}

	// 10. Diff
	// Modify hello.txt in working tree
	if err := os.WriteFile(filepath.Join(repoDir, "hello.txt"), []byte("hello modified\n"), 0644); err != nil {
		t.Fatal(err)
	}
	diffRes, err := backend.Diff(ctx, DiffOptions{})
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}
	if len(diffRes.Files) != 1 || diffRes.Files[0].NewName != "hello.txt" {
		t.Errorf("unexpected diff result: %+v", diffRes)
	}

	// 11. Branch operations
	if err := backend.Checkout(ctx, "feature", CheckoutOptions{Create: true}); err != nil {
		t.Fatalf("Checkout -b failed: %v", err)
	}
	cur, _ := backend.CurrentBranch(ctx)
	if cur != "feature" {
		t.Errorf("expected feature branch, got %q", cur)
	}

	branches, err := backend.Branches(ctx, BranchOptions{})
	if err != nil || len(branches) < 2 {
		t.Errorf("Branches failed: %v, got %+v", err, branches)
	}

	// Switch back to main
	if err := backend.Checkout(ctx, "main", CheckoutOptions{}); err != nil {
		t.Fatalf("Checkout main failed: %v", err)
	}

	// 12. Tags
	if err := backend.CreateTag(ctx, TagOptions{Create: "v1.0.0"}); err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}
	tags, err := backend.Tags(ctx, TagOptions{List: true})
	if err != nil || len(tags) != 1 || tags[0].Name != "v1.0.0" {
		t.Errorf("Tags failed: %v, got %+v", err, tags)
	}
	if err := backend.DeleteTag(ctx, "v1.0.0"); err != nil {
		t.Fatalf("DeleteTag failed: %v", err)
	}

	// 13. Config
	if err := backend.ConfigSet(ctx, "gitforge.testkey", "testvalue", false); err != nil {
		t.Fatalf("ConfigSet failed: %v", err)
	}
	val, err := backend.ConfigGet(ctx, "gitforge.testkey", false)
	if err != nil || val != "testvalue" {
		t.Errorf("ConfigGet failed: %v, got %q", err, val)
	}
	cfgMap, err := backend.ConfigList(ctx, false)
	if err != nil || cfgMap["gitforge.testkey"] != "testvalue" {
		t.Errorf("ConfigList failed: %v, map: %+v", err, cfgMap)
	}

	// 14. Remotes
	if err := backend.RemoteAdd(ctx, "origin", "https://github.com/example/repo.git"); err != nil {
		t.Fatalf("RemoteAdd failed: %v", err)
	}
	remotes, err := backend.Remotes(ctx)
	if err != nil || len(remotes) != 1 || remotes[0].Name != "origin" {
		t.Errorf("Remotes failed: %v, got %+v", err, remotes)
	}
	if err := backend.RemoteRemove(ctx, "origin"); err != nil {
		t.Fatalf("RemoteRemove failed: %v", err)
	}

	// 15. Stash
	if err := os.WriteFile(filepath.Join(repoDir, "stashme.txt"), []byte("stash content\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_ = backend.Add(ctx, []string{"stashme.txt"})
	if _, err := backend.Stash(ctx, StashOptions{Action: StashPush, Message: "test stash"}); err != nil {
		t.Fatalf("StashPush failed: %v", err)
	}
	stashes, err := backend.Stash(ctx, StashOptions{Action: StashList})
	if err != nil || len(stashes) != 1 {
		t.Errorf("StashList failed: %v, got %+v", err, stashes)
	}
	if _, err := backend.Stash(ctx, StashOptions{Action: StashPop}); err != nil {
		t.Fatalf("StashPop failed: %v", err)
	}

	// 16. ApplyPatch
	patch := `diff --git a/hello.txt b/hello.txt
--- a/hello.txt
+++ b/hello.txt
@@ -1,1 +1,1 @@
-hello modified
+hello patched
`
	if err := backend.ApplyPatch(ctx, patch, false); err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}

	// 17. Init
	tempNewDir := filepath.Join(t.TempDir(), "new-sub-repo")
	if err := backend.Init(ctx, tempNewDir, false); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 18. Local bare remote for Push, Fetch, Pull, Clone
	bareDir := filepath.Join(t.TempDir(), "bare-remote.git")
	if err := backend.Init(ctx, bareDir, true); err != nil {
		t.Fatalf("Init bare failed: %v", err)
	}
	if err := backend.RemoteAdd(ctx, "localbare", bareDir); err != nil {
		t.Fatalf("RemoteAdd bare failed: %v", err)
	}
	if err := backend.Push(ctx, "localbare", "main", PushOptions{SetUpstream: true}); err != nil {
		t.Fatalf("Push failed: %v", err)
	}
	if err := backend.Fetch(ctx, "localbare", FetchOptions{}); err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}
	if err := backend.Pull(ctx, "localbare", "main", PullOptions{}); err != nil {
		t.Fatalf("Pull failed: %v", err)
	}
	cloneDir := filepath.Join(t.TempDir(), "cloned-repo")
	if err := backend.Clone(ctx, bareDir, cloneDir); err != nil {
		t.Fatalf("Clone failed: %v", err)
	}

	// 19. Argument injection check: ref starting with - must be rejected
	if err := ValidateRef("-rf"); err == nil {
		t.Error("expected ValidateRef to reject '-rf'")
	}
	if _, err := backend.Show(ctx, "--help"); err == nil {
		t.Error("expected backend.Show to reject '--help'")
	}
}
