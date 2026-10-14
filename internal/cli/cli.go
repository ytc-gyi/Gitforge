// Package cli implements the command-line interface for GitForge.
package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/phantom1785322-hub/gitforge/internal/config"
	"github.com/phantom1785322-hub/gitforge/internal/core"
	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

var (
	// gitBackend is the shared git backend for CLI commands.
	gitBackend git.Backend
	// currentConfig is the loaded configuration.
	currentConfig *config.Config
)

// SetBackend sets the backend used by CLI commands (useful for testing).
func SetBackend(b git.Backend) {
	gitBackend = b
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gitforge",
		Short: "GitForge - A fast, cross-platform, local-first Git client",
		Long: `GitForge is a fast, cross-platform, local-first Git client shipping as a single static binary.
Terminal-first, not terminal-only. Local-first, cloud-optional. Extensible by default.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			currentConfig, err = config.Load()
			if err != nil {
				currentConfig = config.DefaultConfig()
			}

			if gitBackend == nil {
				gitBackend, _ = git.NewCLIBackend("")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return passthroughToGit(cmd.Context(), args)
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Register Phase 1 & Phase 2 commands
	rootCmd.AddCommand(
		newVersionCmd(),
		newDoctorCmd(),
		newStatusCmd(),
		newLogCmd(),
		newDiffCmd(),
		newShowCmd(),
		newConfigCmd(),
		newInitCmd(),
		newCloneCmd(),
		newBranchCmd(),
		newCheckoutCmd(),
		newAddCmd(),
		newCommitCmd(),
		newRemoteCmd(),
		newFetchCmd(),
		newPullCmd(),
		newPushCmd(),
		newStashCmd(),
		newTagCmd(),
		newCompletionCmd(rootCmd),
	)

	return rootCmd
}

// Execute parses arguments, expands aliases, and executes the CLI command.
// If an unknown command is invoked, it passes through to git.
func Execute() error {
	return ExecuteArgs(os.Args[1:])
}

// ExecuteArgs executes with given argument slice.
func ExecuteArgs(args []string) error {
	rootCmd := newRootCmd()

	// Load config early for aliases
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}
	currentConfig = cfg

	// Expand alias if applicable
	if len(args) > 0 {
		subcmd := args[0]
		if expanded, ok := cfg.Aliases[subcmd]; ok {
			expandedParts := strings.Fields(expanded)
			args = append(expandedParts, args[1:]...)
		}
	}

	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		// If command failed or was not found, check if it's an unknown subcommand to pass through to git
		if len(args) > 0 {
			if ptErr := passthroughToGit(context.Background(), args); ptErr == nil {
				return nil
			}
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return nil
}

func passthroughToGit(ctx context.Context, args []string) error {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git binary not found in PATH: %w", err)
	}

	cmd := exec.CommandContext(ctx, gitPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func requireGit() error {
	if gitBackend == nil {
		var err error
		gitBackend, err = git.NewCLIBackend("")
		if err != nil {
			return core.ErrGitNotFound
		}
	}
	return nil
}
