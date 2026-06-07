// Package config defines the "nous config" command.
package config

import (
	"fmt"

	"github.com/cajohnson0125/Nous/internal"
	"github.com/cajohnson0125/Nous/internal/config"

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
			exists, path, err := config.UserConfigExists()
			if err != nil {
				return fmt.Errorf("check config path: %w", err)
			}

			if exists {
				fmt.Printf("Config file already exists: %s\n", path)
				fmt.Println("Edit the existing file or remove it before reinitializing.")
				return nil
			}

			cfg := config.Default()
			written, err := config.Save(cfg)
			if err != nil {
				return fmt.Errorf("write config: %w", err)
			}

			fmt.Printf("Config file created: %s\n", written)
			return nil
		},
	}

	return cmd
}
