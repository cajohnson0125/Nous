// Package chat defines the "nous chat" command.
package chat

import (
	"fmt"

	"github.com/cajohnson0125/Nous/internal"

	"github.com/spf13/cobra"
)

// New creates the chat command.
func New(app *internal.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Open an interactive chat session",
		Long:  "Start an interactive chat session with Nous. Wires to the TUI in Phase 4.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("chat TUI not yet implemented (Phase 4)")
		},
	}

	cmd.GroupID = "core"

	return cmd
}
