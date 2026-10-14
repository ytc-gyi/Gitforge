package main

import (
	"log/slog"
	"os"

	"github.com/phantom1785322-hub/gitforge/internal/cli"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
