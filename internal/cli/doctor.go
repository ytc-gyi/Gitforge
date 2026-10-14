package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/phantom1785322-hub/gitforge/internal/config"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check your environment for potential issues",
		Long:  `Validates Git installation, repository state, configuration, terminal capabilities, and OS environment.`,
		Args:  cobra.NoArgs,
		RunE:  runDoctor,
	}
}

func runDoctor(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	fmt.Println("GitForge Doctor")
	fmt.Println(strings.Repeat("─", 40))

	allOK := true

	// 1. Git check
	allOK = checkGit(ctx) && allOK

	// 2. Config check
	allOK = checkConfig() && allOK

	// 3. Terminal check
	allOK = checkTerminal() && allOK

	// 4. Platform check
	allOK = checkPlatform() && allOK

	fmt.Println()
	if allOK {
		fmt.Println("✅ All checks passed!")
	} else {
		fmt.Println("⚠️  Some checks reported warnings. See details above.")
	}

	return nil
}

func checkGit(ctx context.Context) bool {
	fmt.Println("\n📋 Git")

	gitPath, err := exec.LookPath("git")
	if err != nil {
		fmt.Println("  ❌ git executable not found in PATH")
		fmt.Println("     Fix: Install Git from https://git-scm.com/")
		return false
	}
	fmt.Printf("  ✅ git found: %s\n", gitPath)

	out, err := exec.CommandContext(ctx, gitPath, "version").Output()
	if err != nil {
		fmt.Println("  ❌ Failed to get git version")
		return false
	}
	fmt.Printf("  ✅ %s\n", strings.TrimSpace(string(out)))

	repoCheck := exec.CommandContext(ctx, gitPath, "rev-parse", "--git-dir")
	if err := repoCheck.Run(); err != nil {
		fmt.Println("  ℹ️  Not currently inside a git repository")
	} else {
		fmt.Println("  ✅ Inside a git repository")
	}

	return true
}

func checkConfig() bool {
	fmt.Println("\n📋 Configuration")

	_, err := config.Load()
	if err != nil {
		fmt.Printf("  ❌ Config error: %v\n", err)
		return false
	}

	paths := config.LoadedPaths()
	if len(paths) == 0 {
		fmt.Println("  ℹ️  No config files found (using built-in defaults)")
	} else {
		for _, p := range paths {
			fmt.Printf("  ✅ Loaded: %s\n", p)
		}
	}

	return true
}

func checkTerminal() bool {
	fmt.Println("\n📋 Terminal")

	fi, err := os.Stdout.Stat()
	if err == nil && (fi.Mode()&os.ModeCharDevice) != 0 {
		fmt.Println("  ✅ stdout is a TTY")
	} else {
		fmt.Println("  ℹ️  stdout is not a TTY (piped or redirected)")
	}

	if term := os.Getenv("TERM"); term != "" {
		fmt.Printf("  ✅ TERM=%s\n", term)
	} else {
		fmt.Println("  ℹ️  TERM is not set")
	}

	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		fmt.Println("  ℹ️  NO_COLOR is set (color output disabled)")
	}

	if colorTerm := os.Getenv("COLORTERM"); colorTerm != "" {
		fmt.Printf("  ✅ COLORTERM=%s (true color support)\n", colorTerm)
	}

	return true
}

func checkPlatform() bool {
	fmt.Println("\n📋 Platform")
	fmt.Printf("  ✅ OS: %s\n", runtime.GOOS)
	fmt.Printf("  ✅ Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("  ✅ Go runtime: %s\n", runtime.Version())
	return true
}
