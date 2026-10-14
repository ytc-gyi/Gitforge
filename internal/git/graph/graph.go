package graph

import (
	"strings"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

// Node represents a commit placed in a 2D graph grid.
type Node struct {
	Commit    core.Commit
	Row       int
	Col       int
	Glyph     string
	GraphLine string
	Parents   []ParentLink
}

// ParentLink describes a directed line from this commit to one of its parents.
type ParentLink struct {
	ParentHash string
	ToCol      int
	ToRow      int // -1 if parent is not in the current commit batch
}

// Layout computes the column lane assignment and ASCII rendering for each commit.
// Commits should be provided in topological / chronological order (as output by git log).
func Layout(commits []core.Commit) []Node {
	if len(commits) == 0 {
		return nil
	}

	nodes := make([]Node, len(commits))
	commitPositions := make(map[string][2]int, len(commits)) // hash -> [row, col]

	// activeLanes tracks the commit hash currently expected in each column.
	var activeLanes []string

	for r, commit := range commits {
		// 1. Find or allocate lane for this commit.
		col := -1
		for i, laneHash := range activeLanes {
			if laneHash == commit.Hash {
				col = i
				break
			}
		}

		if col == -1 {
			// Find first empty slot or append.
			for i, laneHash := range activeLanes {
				if laneHash == "" {
					col = i
					activeLanes[i] = commit.Hash
					break
				}
			}
			if col == -1 {
				col = len(activeLanes)
				activeLanes = append(activeLanes, commit.Hash)
			}
		}

		commitPositions[commit.Hash] = [2]int{r, col}

		// 2. Render ASCII glyph & graph line for TUI.
		var sb strings.Builder
		for i := 0; i < len(activeLanes); i++ {
			if i == col {
				sb.WriteString("*")
			} else if activeLanes[i] != "" {
				sb.WriteString("|")
			} else {
				sb.WriteString(" ")
			}
			if i < len(activeLanes)-1 {
				sb.WriteString(" ")
			}
		}

		nodes[r] = Node{
			Commit:    commit,
			Row:       r,
			Col:       col,
			Glyph:     "*",
			GraphLine: sb.String(),
		}

		// 3. Update active lanes for the commit's parents.
		if len(commit.ParentHashes) == 0 {
			// Root commit: free this lane.
			activeLanes[col] = ""
		} else {
			// Primary parent inherits this lane.
			activeLanes[col] = commit.ParentHashes[0]

			// Additional parents (merges) get additional lanes.
			for _, parentHash := range commit.ParentHashes[1:] {
				alreadyAssigned := false
				for _, h := range activeLanes {
					if h == parentHash {
						alreadyAssigned = true
						break
					}
				}
				if !alreadyAssigned {
					// Find an empty lane or append.
					placed := false
					for i, h := range activeLanes {
						if h == "" {
							activeLanes[i] = parentHash
							placed = true
							break
						}
					}
					if !placed {
						activeLanes = append(activeLanes, parentHash)
					}
				}
			}
		}

		// Trim trailing empty lanes.
		for len(activeLanes) > 0 && activeLanes[len(activeLanes)-1] == "" {
			activeLanes = activeLanes[:len(activeLanes)-1]
		}
	}

	// 4. Fill in parent coordinates for SVG / Web graphing.
	for r := range nodes {
		for _, pHash := range nodes[r].Commit.ParentHashes {
			link := ParentLink{
				ParentHash: pHash,
				ToRow:      -1,
				ToCol:      -1,
			}
			if pos, found := commitPositions[pHash]; found {
				link.ToRow = pos[0]
				link.ToCol = pos[1]
			}
			nodes[r].Parents = append(nodes[r].Parents, link)
		}
	}

	return nodes
}
