package git

import (
	"testing"
)

func TestParseBranches(t *testing.T) {
	t.Parallel()

	input := "main\x001234567\x00origin/main\x00ahead 2, behind 1\x00*\n" +
		"feature\x00abcdef0\x00\x00\x00 \n" +
		"bugfix\x007890abc\x00origin/bugfix\x00ahead 3\x00 \n"

	branches, err := ParseBranches(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}

	// main
	b0 := branches[0]
	if b0.Name != "main" || b0.Hash != "1234567" || b0.Upstream != "origin/main" || !b0.Current {
		t.Errorf("branch 0 mismatch: %+v", b0)
	}
	if b0.Ahead != 2 || b0.Behind != 1 {
		t.Errorf("branch 0 ahead/behind mismatch: ahead=%d behind=%d", b0.Ahead, b0.Behind)
	}

	// feature
	b1 := branches[1]
	if b1.Name != "feature" || b1.Hash != "abcdef0" || b1.Upstream != "" || b1.Current {
		t.Errorf("branch 1 mismatch: %+v", b1)
	}

	// bugfix
	b2 := branches[2]
	if b2.Name != "bugfix" || b2.Ahead != 3 || b2.Behind != 0 {
		t.Errorf("branch 2 mismatch: %+v", b2)
	}
}

func TestParseBranchesEmpty(t *testing.T) {
	t.Parallel()

	branches, err := ParseBranches("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 0 {
		t.Errorf("expected 0 branches, got %d", len(branches))
	}
}
