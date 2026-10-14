package graph

import (
	"testing"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

func TestEmptyGraph(t *testing.T) {
	t.Parallel()
	nodes := Layout(nil)
	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(nodes))
	}
}

func TestLinearHistory(t *testing.T) {
	t.Parallel()

	// C3 -> C2 -> C1
	commits := []core.Commit{
		{Hash: "c3", ParentHashes: []string{"c2"}, Subject: "commit 3"},
		{Hash: "c2", ParentHashes: []string{"c1"}, Subject: "commit 2"},
		{Hash: "c1", ParentHashes: nil, Subject: "root commit"},
	}

	nodes := Layout(commits)
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}

	for i, node := range nodes {
		if node.Col != 0 {
			t.Errorf("node %d: expected Col 0, got %d", i, node.Col)
		}
		if node.Row != i {
			t.Errorf("node %d: expected Row %d, got %d", i, i, node.Row)
		}
		if node.GraphLine != "*" {
			t.Errorf("node %d: expected GraphLine '*', got %q", i, node.GraphLine)
		}
	}

	// Check parent links
	if len(nodes[0].Parents) != 1 || nodes[0].Parents[0].ToRow != 1 || nodes[0].Parents[0].ToCol != 0 {
		t.Errorf("node 0 parent link mismatch: %+v", nodes[0].Parents)
	}
	if len(nodes[1].Parents) != 1 || nodes[1].Parents[0].ToRow != 2 || nodes[1].Parents[0].ToCol != 0 {
		t.Errorf("node 1 parent link mismatch: %+v", nodes[1].Parents)
	}
	if len(nodes[2].Parents) != 0 {
		t.Errorf("root node should have 0 parents, got %d", len(nodes[2].Parents))
	}
}

func TestBranchAndMerge(t *testing.T) {
	t.Parallel()

	// Merge commit M (parents: B, F)
	// Feature commit F (parent: B)
	// Base commit B (parent: root)
	// Root commit R
	commits := []core.Commit{
		{Hash: "M", ParentHashes: []string{"B", "F"}, Subject: "Merge branch 'feature'"},
		{Hash: "F", ParentHashes: []string{"B"}, Subject: "Feature work"},
		{Hash: "B", ParentHashes: []string{"R"}, Subject: "Base commit"},
		{Hash: "R", ParentHashes: nil, Subject: "Root commit"},
	}

	nodes := Layout(commits)
	if len(nodes) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(nodes))
	}

	// M is in col 0
	if nodes[0].Col != 0 {
		t.Errorf("M: expected Col 0, got %d", nodes[0].Col)
	}
	// F should be in col 1 (second parent of M)
	if nodes[1].Col != 1 {
		t.Errorf("F: expected Col 1, got %d", nodes[1].Col)
	}
	// B should be in col 0 (first parent of M)
	if nodes[2].Col != 0 {
		t.Errorf("B: expected Col 0, got %d", nodes[2].Col)
	}
	// R should be in col 0
	if nodes[3].Col != 0 {
		t.Errorf("R: expected Col 0, got %d", nodes[3].Col)
	}

	// Check M parents: B (row 2, col 0) and F (row 1, col 1)
	if len(nodes[0].Parents) != 2 {
		t.Fatalf("M expected 2 parents, got %d", len(nodes[0].Parents))
	}
	if nodes[0].Parents[0].ToRow != 2 || nodes[0].Parents[0].ToCol != 0 {
		t.Errorf("M parent 0 link mismatch: %+v", nodes[0].Parents[0])
	}
	if nodes[0].Parents[1].ToRow != 1 || nodes[0].Parents[1].ToCol != 1 {
		t.Errorf("M parent 1 link mismatch: %+v", nodes[0].Parents[1])
	}
}
