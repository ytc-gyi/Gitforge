package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newBranchCmd() *cobra.Command {
	var (
		all        bool
		remotes    bool
		deleteName string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:     "branch [name]",
		Short:   "List, create, or delete branches",
		Aliases: []string{"br"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()

			// Delete branch
			if deleteName != "" {
				opts := git.BranchOptions{Delete: deleteName}
				if _, err := gitBackend.Branches(ctx, opts); err != nil {
					return err
				}
				fmt.Printf("Deleted branch %s\n", FormatBranch(deleteName))
				return nil
			}

			// Create branch
			if len(args) == 1 {
				branchName := args[0]
				opts := git.BranchOptions{Create: branchName}
				if _, err := gitBackend.Branches(ctx, opts); err != nil {
					return err
				}
				fmt.Printf("Created branch %s\n", FormatBranch(branchName))
				return nil
			}

			// List branches
			opts := git.BranchOptions{
				All:     all,
				Remotes: remotes,
			}
			branches, err := gitBackend.Branches(ctx, opts)
			if err != nil {
				return err
			}

			if jsonOutput {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(branches)
			}

			for _, b := range branches {
				prefix := "  "
				name := b.Name
				if b.Current {
					prefix = "* "
					if IsTTY() {
						name = StyleSuccess.Render(b.Name)
					}
				}

				trackInfo := ""
				if b.Upstream != "" {
					trackInfo = fmt.Sprintf(" [%s]", b.Upstream)
					if b.Ahead > 0 || b.Behind > 0 {
						trackInfo += fmt.Sprintf(" (ahead %d, behind %d)", b.Ahead, b.Behind)
					}
				}

				fmt.Printf("%s%s %s%s\n", prefix, FormatHash(b.Hash), name, StyleMuted.Render(trackInfo))
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "List both remote-tracking and local branches")
	cmd.Flags().BoolVarP(&remotes, "remotes", "r", false, "List remote-tracking branches")
	cmd.Flags().StringVarP(&deleteName, "delete", "d", "", "Delete a branch")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output branches as JSON")

	return cmd
}
