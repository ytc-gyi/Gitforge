package cli

import (
	"fmt"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newCommitCmd() *cobra.Command {
	var (
		message string
		amend   bool
	)

	cmd := &cobra.Command{
		Use:     "commit [-m <message>] [--amend]",
		Short:   "Record changes to the repository",
		Aliases: []string{"cm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			opts := git.CommitOptions{
				Message: message,
				Amend:   amend,
			}
			hash, err := gitBackend.Commit(ctx, opts)
			if err != nil {
				return err
			}
			if amend {
				fmt.Printf("[%s] %s\n", FormatHash(hash), FormatSuccess("Amended commit"))
			} else {
				fmt.Printf("[%s] %s\n", FormatHash(hash), message)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Commit message")
	cmd.Flags().BoolVar(&amend, "amend", false, "Amend the previous commit")

	return cmd
}
