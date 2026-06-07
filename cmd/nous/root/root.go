// Package root defines the root Cobra command for Nous.
package root

import (
	"context"

	"github.com/cajohnson0125/Nous/internal"
	"github.com/cajohnson0125/Nous/cmd/nous/chat"
	"github.com/cajohnson0125/Nous/cmd/nous/config"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

// New creates the root command with all subcommands registered.
func New(app *internal.App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nous",
		Short: "Nous — multi-agent state engine",
		Long:  "Nous is a multi-layer, dynamic state engine for terminal workflows.",
	}

	rootCmd.AddGroup(
		&cobra.Group{ID: "core", Title: "Core Commands"},
	)

	rootCmd.AddCommand(chat.New(app), config.New(app))

	return rootCmd
}

// Execute runs the root command via Fang for styled output.
func Execute(app *internal.App) error {
	rootCmd := New(app)

	return fang.Execute(context.Background(), rootCmd,
		fang.WithVersion(app.Version),
		fang.WithCommit(app.Commit),
	)
}
