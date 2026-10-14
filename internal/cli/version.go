package cli

import (
	"fmt"

	"github.com/phantom1785322-hub/gitforge/internal/core"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the GitForge version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(core.VersionString())
			return nil
		},
	}
}
