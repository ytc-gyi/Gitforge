package cli

import (
	"fmt"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newPullCmd() *cobra.Command {
	var rebase bool

	cmd := &cobra.Command{
		Use:   "pull [remote] [branch]",
		Short: "Fetch from and integrate with another repository or a local branch",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			remote := ""
			branch := ""
			if len(args) > 0 {
				remote = args[0]
			}
			if len(args) > 1 {
				branch = args[1]
			}
			opts := git.PullOptions{Rebase: rebase}
			if err := gitBackend.Pull(ctx, remote, branch, opts); err != nil {
				return err
			}
			fmt.Printf("%s Pull completed\n", FormatSuccess("✅"))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&rebase, "rebase", "r", false, "Rebase current branch on top of upstream branch after fetch")

	return cmd
}
