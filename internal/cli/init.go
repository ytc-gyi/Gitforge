package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var bare bool

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize a new Git repository",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			ctx := cmd.Context()
			if err := gitBackend.Init(ctx, path, bare); err != nil {
				return err
			}
			if bare {
				fmt.Printf("%s Initialized empty bare Git repository in %s\n", FormatSuccess("✅"), path)
			} else {
				fmt.Printf("%s Initialized empty Git repository in %s\n", FormatSuccess("✅"), path)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&bare, "bare", false, "Create a bare repository")

	return cmd
}
