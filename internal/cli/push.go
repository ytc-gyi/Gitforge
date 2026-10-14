package cli

import (
	"fmt"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newPushCmd() *cobra.Command {
	var (
		force       bool
		setUpstream bool
		tags        bool
	)

	cmd := &cobra.Command{
		Use:   "push [remote] [branch]",
		Short: "Update remote refs along with associated objects",
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
			opts := git.PushOptions{
				Force:       force,
				SetUpstream: setUpstream,
				Tags:        tags,
			}
			if err := gitBackend.Push(ctx, remote, branch, opts); err != nil {
				return err
			}
			fmt.Printf("%s Push completed\n", FormatSuccess("✅"))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force update of remote refs")
	cmd.Flags().BoolVarP(&setUpstream, "set-upstream", "u", false, "Add upstream (tracking) reference")
	cmd.Flags().BoolVar(&tags, "tags", false, "Push all tags")

	return cmd
}
