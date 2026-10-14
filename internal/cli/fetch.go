package cli

import (
	"fmt"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newFetchCmd() *cobra.Command {
	var (
		all   bool
		prune bool
		tags  bool
	)

	cmd := &cobra.Command{
		Use:   "fetch [remote]",
		Short: "Download objects and refs from another repository",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			remote := ""
			if len(args) > 0 {
				remote = args[0]
			}
			opts := git.FetchOptions{
				All:   all,
				Prune: prune,
				Tags:  tags,
			}
			if err := gitBackend.Fetch(ctx, remote, opts); err != nil {
				return err
			}
			fmt.Printf("%s Fetch completed\n", FormatSuccess("✅"))
			return nil
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Fetch all remotes")
	cmd.Flags().BoolVarP(&prune, "prune", "p", false, "Before fetching, remove any remote-tracking references that no longer exist on the remote")
	cmd.Flags().BoolVarP(&tags, "tags", "t", false, "Fetch all tags from the remote")

	return cmd
}
