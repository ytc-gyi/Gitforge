package cli

import (
	"fmt"

	"github.com/phantom1785322-hub/gitforge/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	var global bool

	cmd := &cobra.Command{
		Use:   "config [key] [value]",
		Short: "View or modify configuration",
		Long: `View or modify GitForge and git configuration.
Without arguments, shows all loaded GitForge configuration (sensitive secrets are redacted).
With one argument, gets the configuration value for that key.
With two arguments, sets the value for that key.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				return showAllConfig()
			case 1:
				if err := requireGit(); err != nil {
					return err
				}
				return getConfigValue(cmd, args[0], global)
			case 2:
				if err := requireGit(); err != nil {
					return err
				}
				return setConfigValue(cmd, args[0], args[1], global)
			default:
				return fmt.Errorf("expected 0 to 2 arguments, got %d", len(args))
			}
		},
	}

	cmd.Flags().BoolVar(&global, "global", false, "Use global git config")

	return cmd
}

func showAllConfig() error {
	cfg := config.Get()
	redacted := config.RedactedConfig(cfg)

	fmt.Println("GitForge Configuration (secrets redacted):")
	fmt.Println()

	// [core]
	fmt.Println("[core]")
	fmt.Printf("  editor = %q\n", redacted.Core.Editor)
	fmt.Printf("  pager = %q\n", redacted.Core.Pager)
	fmt.Printf("  auto_fetch = %v\n", redacted.Core.AutoFetch)
	fmt.Printf("  auto_fetch_interval = %d\n", redacted.Core.AutoFetchInterval)
	fmt.Println()

	// [ui]
	fmt.Println("[ui]")
	fmt.Printf("  theme = %q\n", redacted.UI.Theme)
	fmt.Printf("  icons = %q\n", redacted.UI.Icons)
	fmt.Printf("  mouse = %v\n", redacted.UI.Mouse)
	fmt.Printf("  unicode = %v\n", redacted.UI.Unicode)
	fmt.Printf("  compact_mode = %v\n", redacted.UI.CompactMode)
	fmt.Println()

	// [ai]
	fmt.Println("[ai]")
	fmt.Printf("  enabled = %v\n", redacted.AI.Enabled)
	fmt.Printf("  engine = %q\n", redacted.AI.Engine)
	fmt.Printf("  endpoint = %q\n", redacted.AI.Endpoint)
	fmt.Printf("  model = %q\n", redacted.AI.Model)
	fmt.Printf("  context_lines = %d\n", redacted.AI.ContextLines)
	fmt.Println()

	// [github]
	fmt.Println("[github]")
	fmt.Printf("  token = %q\n", redacted.GitHub.Token)
	fmt.Printf("  host = %q\n", redacted.GitHub.Host)
	fmt.Println()

	// [gitlab]
	fmt.Println("[gitlab]")
	fmt.Printf("  token = %q\n", redacted.GitLab.Token)
	fmt.Printf("  host = %q\n", redacted.GitLab.Host)
	fmt.Println()

	// [plugins]
	fmt.Println("[plugins]")
	fmt.Printf("  directory = %q\n", redacted.Plugins.Directory)
	fmt.Printf("  allow_unsigned = %v\n", redacted.Plugins.AllowUnsigned)
	fmt.Println()

	// [aliases]
	if len(redacted.Aliases) > 0 {
		fmt.Println("[aliases]")
		for k, v := range redacted.Aliases {
			fmt.Printf("  %s = %q\n", k, v)
		}
		fmt.Println()
	}

	// Loaded config paths
	paths := config.LoadedPaths()
	if len(paths) > 0 {
		fmt.Println("Loaded files:")
		for _, p := range paths {
			fmt.Printf("  - %s\n", p)
		}
	}

	return nil
}

func getConfigValue(cmd *cobra.Command, key string, global bool) error {
	ctx := cmd.Context()
	val, err := gitBackend.ConfigGet(ctx, key, global)
	if err != nil {
		return fmt.Errorf("getting config %q: %w", key, err)
	}

	val = config.RedactValue(key, val)
	fmt.Println(val)
	return nil
}

func setConfigValue(cmd *cobra.Command, key, value string, global bool) error {
	ctx := cmd.Context()
	if err := gitBackend.ConfigSet(ctx, key, value, global); err != nil {
		return fmt.Errorf("setting config %q: %w", key, err)
	}
	return nil
}
