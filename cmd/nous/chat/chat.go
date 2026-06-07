// Package chat defines the "nous chat" command.
package chat

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"github.com/cajohnson0125/Nous/internal"
	"github.com/cajohnson0125/Nous/internal/config"
	"github.com/cajohnson0125/Nous/internal/layout"

	"github.com/spf13/cobra"
)

// New creates the chat command.
func New(app *internal.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Open an interactive chat session",
		Long:  "Start an interactive chat session with Nous. Type a message and press Enter to send.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			m := layout.NewModel(cfg)
			p := tea.NewProgram(m)

			_, err = p.Run()
			return err
		},
	}

	cmd.GroupID = "core"

	return cmd
}
