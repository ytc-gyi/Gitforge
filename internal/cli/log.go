package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/phantom1785322-hub/gitforge/internal/core"
	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/phantom1785322-hub/gitforge/internal/git/graph"
	"github.com/spf13/cobra"
)

func newLogCmd() *cobra.Command {
	var (
		maxCount   int
		graphMode  bool
		oneline    bool
		jsonOutput bool
		all        bool
	)

	cmd := &cobra.Command{
		Use:     "log [ref] [-- paths...]",
		Short:   "Show commit logs",
		Aliases: []string{"lg"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			opts := git.LogOptions{
				MaxCount: maxCount,
				All:      all,
				Oneline:  oneline,
			}
			for i, arg := range args {
				if arg == "--" {
					if i > 0 {
						opts.Ref = args[0]
					}
					if i+1 < len(args) {
						opts.Paths = args[i+1:]
					}
					break
				}
				if i == 0 {
					opts.Ref = arg
				}
			}

			return runLog(cmd.Context(), opts, graphMode, oneline, jsonOutput)
		},
	}

	cmd.Flags().IntVarP(&maxCount, "max-count", "n", 20, "Limit the number of commits")
	cmd.Flags().BoolVar(&graphMode, "graph", false, "Draw a text-based graphical representation")
	cmd.Flags().BoolVar(&oneline, "oneline", false, "Show each commit on a single line")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output log as JSON")
	cmd.Flags().BoolVar(&all, "all", false, "Show all commits from all branches")

	return cmd
}

func runLog(ctx context.Context, opts git.LogOptions, graphMode, oneline, jsonOut bool) error {
	commits, err := gitBackend.Log(ctx, opts)
	if err != nil {
		return fmt.Errorf("getting log: %w", err)
	}

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(commits)
	}

	if graphMode {
		nodes := graph.Layout(commits)
		for _, node := range nodes {
			c := node.Commit
			refs := formatRefs(c.Refs)
			if refs != "" {
				fmt.Printf("%s %s %s %s\n", node.GraphLine, c.AbbrevHash, refs, c.Subject)
			} else {
				fmt.Printf("%s %s %s\n", node.GraphLine, c.AbbrevHash, c.Subject)
			}
		}
		return nil
	}

	for _, c := range commits {
		if oneline {
			refs := formatRefs(c.Refs)
			if refs != "" {
				fmt.Printf("%s %s %s\n", c.AbbrevHash, refs, c.Subject)
			} else {
				fmt.Printf("%s %s\n", c.AbbrevHash, c.Subject)
			}
		} else {
			fmt.Printf("commit %s\n", c.Hash)
			if len(c.Refs) > 0 {
				fmt.Printf("Refs: %s\n", formatRefs(c.Refs))
			}
			fmt.Printf("Author: %s <%s>\n", c.Author.Name, c.Author.Email)
			fmt.Printf("Date:   %s\n", c.Author.Date.Format("Mon Jan 2 15:04:05 2006 -0700"))
			fmt.Println()
			for _, line := range strings.Split(c.Subject, "\n") {
				fmt.Printf("    %s\n", line)
			}
			if c.Body != "" {
				fmt.Println()
				for _, line := range strings.Split(c.Body, "\n") {
					fmt.Printf("    %s\n", line)
				}
			}
			fmt.Println()
		}
	}

	return nil
}

func formatRefs(refs []core.Ref) string {
	if len(refs) == 0 {
		return ""
	}
	var parts []string
	for _, r := range refs {
		switch r.Type {
		case core.RefTypeHEAD:
			if r.Name != "HEAD" {
				parts = append(parts, fmt.Sprintf("HEAD -> %s", r.Name))
			} else {
				parts = append(parts, "HEAD")
			}
		case core.RefTypeTag:
			parts = append(parts, fmt.Sprintf("tag: %s", r.Name))
		default:
			parts = append(parts, r.Name)
		}
	}
	return "(" + strings.Join(parts, ", ") + ")"
}
