package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCloneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone <repository> [directory]",
		Short: "Clone a repository into a new directory",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			url := args[0]
			path := ""
			if len(args) > 1 {
				path = args[1]
			}
			ctx := cmd.Context()
			fmt.Printf("Cloning into %q...\n", url)
			if err := gitBackend.Clone(ctx, url, path); err != nil {
				return err
			}
			fmt.Printf("%s Successfully cloned repository\n", FormatSuccess("✅"))
			return nil
		},
	}
}
