package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "show [ref]",
		Short: "Show a commit and its changes",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ref := "HEAD"
			if len(args) > 0 {
				ref = args[0]
			}

			if err := git.ValidateRef(ref); err != nil {
				return err
			}

			ctx := cmd.Context()
			detail, err := gitBackend.Show(ctx, ref)
			if err != nil {
				return fmt.Errorf("showing commit: %w", err)
			}

			if jsonOutput {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(detail)
			}

			fmt.Printf("commit %s\n", detail.Hash)
			fmt.Printf("Author: %s <%s>\n", detail.Author.Name, detail.Author.Email)
			fmt.Printf("Date:   %s\n", detail.Author.Date.Format("Mon Jan 2 15:04:05 2006 -0700"))
			fmt.Println()
			for _, line := range strings.Split(detail.Subject, "\n") {
				fmt.Printf("    %s\n", line)
			}
			if detail.Body != "" {
				fmt.Println()
				for _, line := range strings.Split(detail.Body, "\n") {
					fmt.Printf("    %s\n", line)
				}
			}
			fmt.Println()

			// Show raw diff
			rawDiff, err := gitBackend.DiffRaw(ctx, git.DiffOptions{Ref: ref + "~1", RefTo: ref})
			if err == nil && rawDiff != "" {
				fmt.Print(rawDiff)
			} else {
				// Initial commit (no parent), show diff for this commit
				rawDiff, err = gitBackend.DiffRaw(ctx, git.DiffOptions{Ref: ref})
				if err == nil && rawDiff != "" {
					fmt.Print(rawDiff)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return cmd
}
