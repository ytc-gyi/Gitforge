package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/phantom1785322-hub/gitforge/internal/core"
	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Show the working tree status",
		Aliases: []string{"st"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			return runStatus(cmd.Context(), jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output status as JSON")

	return cmd
}

func runStatus(ctx context.Context, jsonOut bool) error {
	entries, err := gitBackend.Status(ctx, git.StatusOptions{})
	if err != nil {
		return fmt.Errorf("getting status: %w", err)
	}

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	if len(entries) == 0 {
		fmt.Println("nothing to commit, working tree clean")
		return nil
	}

	var staged, unstaged, untracked []core.StatusEntry
	for _, e := range entries {
		switch {
		case e.Index == core.StatusUntracked:
			untracked = append(untracked, e)
		case e.IsStaged():
			staged = append(staged, e)
			if e.IsUnstaged() {
				unstaged = append(unstaged, e)
			}
		case e.IsUnstaged():
			unstaged = append(unstaged, e)
		}
	}

	if len(staged) > 0 {
		fmt.Println("Changes to be committed:")
		for _, e := range staged {
			fmt.Printf("  %s %s\n", statusLabel(e.Index), formatEntryPath(e))
		}
		fmt.Println()
	}

	if len(unstaged) > 0 {
		fmt.Println("Changes not staged for commit:")
		for _, e := range unstaged {
			fmt.Printf("  %s %s\n", statusLabel(e.Worktree), e.Path)
		}
		fmt.Println()
	}

	if len(untracked) > 0 {
		fmt.Println("Untracked files:")
		for _, e := range untracked {
			fmt.Printf("  %s\n", e.Path)
		}
		fmt.Println()
	}

	return nil
}

func formatEntryPath(e core.StatusEntry) string {
	if e.OrigPath != "" {
		return fmt.Sprintf("%s -> %s", e.OrigPath, e.Path)
	}
	return e.Path
}

func statusLabel(code core.StatusCode) string {
	switch code {
	case core.StatusModified:
		return "modified:   "
	case core.StatusAdded:
		return "new file:   "
	case core.StatusDeleted:
		return "deleted:    "
	case core.StatusRenamed:
		return "renamed:    "
	case core.StatusCopied:
		return "copied:     "
	case core.StatusTypeChanged:
		return "typechange: "
	default:
		return string(code) + ":          "
	}
}
