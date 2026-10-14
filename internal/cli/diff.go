package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newDiffCmd() *cobra.Command {
	var (
		cached     bool
		jsonOutput bool
		stat       bool
		nameOnly   bool
	)

	cmd := &cobra.Command{
		Use:   "diff [ref] [-- paths...]",
		Short: "Show changes between commits, commit and working tree, etc.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			opts := git.DiffOptions{
				Cached:   cached,
				Stat:     stat,
				NameOnly: nameOnly,
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
				if i == len(args)-1 {
					opts.Ref = args[0]
				}
			}

			return runDiff(cmd.Context(), opts, jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&cached, "cached", false, "View staged changes")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output parsed diff as JSON")
	cmd.Flags().BoolVar(&stat, "stat", false, "Generate a diffstat summary")
	cmd.Flags().BoolVar(&nameOnly, "name-only", false, "Show only names of changed files")

	return cmd
}

func runDiff(ctx context.Context, opts git.DiffOptions, jsonOut bool) error {
	if jsonOut {
		result, err := gitBackend.Diff(ctx, opts)
		if err != nil {
			return fmt.Errorf("getting diff: %w", err)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	raw, err := gitBackend.DiffRaw(ctx, opts)
	if err != nil {
		return fmt.Errorf("getting diff: %w", err)
	}

	if raw != "" {
		fmt.Print(raw)
	}
	return nil
}
