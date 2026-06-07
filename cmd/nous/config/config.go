// Package config defines the "nous config" command.
package config

import (
	"fmt"

	"github.com/cajohnson0125/Nous/internal"

	"github.com/spf13/cobra"
)

// New creates the config command.
func New(app *internal.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Nous configuration",
		Long:  "Commands for initializing and managing the Nous configuration file.",
	}

	cmd.AddCommand(configInitCmd())

	return cmd
}

func configInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a default config file",
		Long:  "Create a default configuration file at the user-level XDG config path.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("config init not yet implemented (Phase 2)")
		},
	}

	return cmd
}
