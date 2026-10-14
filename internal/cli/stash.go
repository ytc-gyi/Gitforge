package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newStashCmd() *cobra.Command {
	var (
		message    string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "stash [command]",
		Short: "Stash the changes in a dirty working directory away",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			// Default action when no subcommand: stash push
			opts := git.StashOptions{
				Action:  git.StashPush,
				Message: message,
			}
			_, err := gitBackend.Stash(ctx, opts)
			if err != nil {
				return err
			}
			fmt.Println("Saved working directory and index state")
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Stash description message")

	// Subcommands
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List the stash entries that you currently have",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			stashes, err := gitBackend.Stash(ctx, git.StashOptions{Action: git.StashList})
			if err != nil {
				return err
			}

			if jsonOutput {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(stashes)
			}

			for _, s := range stashes {
				fmt.Printf("stash@{%d}: %s\n", s.Index, s.Message)
			}
			return nil
		},
	}
	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output stash list as JSON")

	popCmd := &cobra.Command{
		Use:   "pop",
		Short: "Remove a single stashed state from the stash list and apply it on top of the current working tree",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			_, err := gitBackend.Stash(ctx, git.StashOptions{Action: git.StashPop})
			if err != nil {
				return err
			}
			fmt.Println("Applied and dropped stash")
			return nil
		},
	}

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Like pop, but do not remove the state from the stash list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			_, err := gitBackend.Stash(ctx, git.StashOptions{Action: git.StashApply})
			if err != nil {
				return err
			}
			fmt.Println("Applied stash")
			return nil
		},
	}

	dropCmd := &cobra.Command{
		Use:   "drop",
		Short: "Remove a single stash entry from the list of stash entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()
			_, err := gitBackend.Stash(ctx, git.StashOptions{Action: git.StashDrop})
			if err != nil {
				return err
			}
			fmt.Println("Dropped stash")
			return nil
		},
	}

	cmd.AddCommand(listCmd, popCmd, applyCmd, dropCmd)

	return cmd
}
