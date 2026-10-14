package cli

import (
	"fmt"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newCheckoutCmd() *cobra.Command {
	var create bool

	cmd := &cobra.Command{
		Use:     "checkout [-b] <branch/commit/tag>",
		Short:   "Switch branches or restore working tree files",
		Aliases: []string{"co"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ref := args[0]
			ctx := cmd.Context()
			opts := git.CheckoutOptions{Create: create}
			if err := gitBackend.Checkout(ctx, ref, opts); err != nil {
				return err
			}
			if create {
				fmt.Printf("Switched to a new branch %s\n", FormatBranch(ref))
			} else {
				fmt.Printf("Switched to branch %s\n", FormatBranch(ref))
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&create, "branch", "b", false, "Create and checkout a new branch")

	return cmd
}
