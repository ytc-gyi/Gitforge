package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/phantom1785322-hub/gitforge/internal/core"
	"github.com/phantom1785322-hub/gitforge/internal/git"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var patchMode bool

	cmd := &cobra.Command{
		Use:   "add [paths...]",
		Short: "Add file contents to the index",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireGit(); err != nil {
				return err
			}
			ctx := cmd.Context()

			if patchMode {
				return runInteractiveAddPatch(ctx)
			}

			if len(args) == 0 {
				args = []string{"."}
			}

			return gitBackend.Add(ctx, args)
		},
	}

	cmd.Flags().BoolVarP(&patchMode, "patch", "p", false, "Interactively choose hunks of patch between the index and the work tree")

	return cmd
}

func runInteractiveAddPatch(ctx context.Context) error {
	diffRes, err := gitBackend.Diff(ctx, git.DiffOptions{})
	if err != nil {
		return fmt.Errorf("getting diff for interactive staging: %w", err)
	}

	if len(diffRes.Files) == 0 {
		fmt.Println("No unstaged changes to stage.")
		return nil
	}

	scanner := bufio.NewScanner(os.Stdin)

	for _, file := range diffRes.Files {
		if file.Binary {
			fmt.Printf("Skipping binary file %s\n", file.NewName)
			continue
		}

		stageAllRemaining := false
		discardAllRemaining := false

		for i, hunk := range file.Hunks {
			if discardAllRemaining {
				break
			}
			if stageAllRemaining {
				patch := GenerateHunkPatch(file, hunk)
				_ = gitBackend.ApplyPatch(ctx, patch, true)
				continue
			}

			// Display hunk
			fmt.Println(StyleBold.Render(fmt.Sprintf("diff --git a/%s b/%s", file.OldName, file.NewName)))
			fmt.Println(StyleMuted.Render(fmt.Sprintf("@@ -%d,%d +%d,%d @@ (hunk %d of %d)", hunk.OldStart, hunk.OldLines, hunk.NewStart, hunk.NewLines, i+1, len(file.Hunks))))
			for _, line := range hunk.Lines {
				switch line.Type {
				case core.DiffLineAdded:
					fmt.Println(StyleAdded.Render("+" + line.Content))
				case core.DiffLineDeleted:
					fmt.Println(StyleDeleted.Render("-" + line.Content))
				case core.DiffLineContext:
					fmt.Println(" " + line.Content)
				}
			}

			for {
				fmt.Print("Stage this hunk [y,n,q,a,d,?]? ")
				if !scanner.Scan() {
					return nil
				}
				ans := strings.TrimSpace(strings.ToLower(scanner.Text()))
				if ans == "" {
					continue
				}

				switch ans {
				case "y":
					patch := GenerateHunkPatch(file, hunk)
					if err := gitBackend.ApplyPatch(ctx, patch, true); err != nil {
						fmt.Printf("%s Failed to apply patch: %v\n", FormatError("❌"), err)
					} else {
						fmt.Println(FormatSuccess("Hunk staged."))
					}
					goto nextHunk
				case "n":
					goto nextHunk
				case "q":
					return nil
				case "a":
					stageAllRemaining = true
					patch := GenerateHunkPatch(file, hunk)
					_ = gitBackend.ApplyPatch(ctx, patch, true)
					goto nextHunk
				case "d":
					discardAllRemaining = true
					goto nextHunk
				case "?":
					fmt.Println("y - stage this hunk")
					fmt.Println("n - do not stage this hunk")
					fmt.Println("q - quit; do not stage this hunk or any of the remaining ones")
					fmt.Println("a - stage this hunk and all later hunks in the file")
					fmt.Println("d - do not stage this hunk or any of the later hunks in the file")
					fmt.Println("? - print help")
				default:
					fmt.Println("Invalid response. Type '?' for help.")
				}
			}

		nextHunk:
		}
	}

	return nil
}

// GenerateHunkPatch builds a valid unified diff patch string for a single hunk.
func GenerateHunkPatch(file core.DiffFile, hunk core.DiffHunk) string {
	var sb strings.Builder
	oldName := file.OldName
	if oldName == "" {
		oldName = file.NewName
	}
	newName := file.NewName
	if newName == "" {
		newName = file.OldName
	}

	sb.WriteString(fmt.Sprintf("--- a/%s\n", oldName))
	sb.WriteString(fmt.Sprintf("+++ b/%s\n", newName))
	sb.WriteString(hunk.Header + "\n")
	for _, l := range hunk.Lines {
		switch l.Type {
		case core.DiffLineAdded:
			sb.WriteString("+" + l.Content + "\n")
		case core.DiffLineDeleted:
			sb.WriteString("-" + l.Content + "\n")
		case core.DiffLineContext:
			sb.WriteString(" " + l.Content + "\n")
		}
	}
	return sb.String()
}
