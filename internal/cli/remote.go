package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRemoteCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage set of tracked repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			remotes, err := gitBackend.Remotes(ctx)
			if err != nil {
				return err
			}

			if jsonOutput {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(remotes)
			}

			for _, r := range remotes {
				fmt.Printf("%s\t%s (fetch)\n", r.Name, r.FetchURL)
				fmt.Printf("%s\t%s (push)\n", r.Name, r.PushURL)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output remotes as JSON")

	// Subcommands: add and remove
	addCmd := &cobra.Command{
		Use:   "add <name> <url>",
		Short: "Add a new remote",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			name := args[0]
			url := args[1]
			if err := gitBackend.RemoteAdd(ctx, name, url); err != nil {
				return err
			}
			fmt.Printf("Added remote %s (%s)\n", FormatBranch(name), url)
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:     "remove <name>",
		Short:   "Remove a remote",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			name := args[0]
			if err := gitBackend.RemoteRemove(ctx, name); err != nil {
				return err
			}
			fmt.Printf("Removed remote %s\n", FormatBranch(name))
			return nil
		},
	}

	cmd.AddCommand(addCmd, removeCmd)

	return cmd
}
