package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newTagCmd() *cobra.Command {
	var (
		deleteName string
		annotated  bool
		message    string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "tag [name]",
		Short: "Create, list, or delete a tag object",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()

			// Delete tag
			if deleteName != "" {
				if err := gitBackend.DeleteTag(ctx, deleteName); err != nil {
					return err
				}
				fmt.Printf("Deleted tag %s\n", deleteName)
				return nil
			}

			// Create tag
			if len(args) == 1 {
				tagName := args[0]
				opts := git.TagOptions{
					Create:    tagName,
					Annotated: annotated || message != "",
					Message:   message,
				}
				if err := gitBackend.CreateTag(ctx, opts); err != nil {
					return err
				}
				fmt.Printf("Created tag %s\n", tagName)
				return nil
			}

			// List tags
			tags, err := gitBackend.Tags(ctx, git.TagOptions{List: true})
			if err != nil {
				return err
			}

			if jsonOutput {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(tags)
			}

			for _, t := range tags {
				if t.IsAnnotated && t.Message != "" {
					fmt.Printf("%s %s %s\n", FormatHash(t.Hash), t.Name, StyleMuted.Render(t.Message))
				} else {
					fmt.Printf("%s %s\n", FormatHash(t.Hash), t.Name)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&deleteName, "delete", "d", "", "Delete a tag")
	cmd.Flags().BoolVarP(&annotated, "annotate", "a", false, "Make an unsigned, annotated tag object")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Tag message for annotated tag")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output tags as JSON")

	return cmd
}
