package styles

import "charm.land/lipgloss/v2"

// Base styles used throughout the Nous TUI.
//
// Styles are immutable value types (Lip Gloss default). Consumers copy
// and override as needed — never mutate these directly.
var (
	// Title styles the "Nous" header shown in the viewport.
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary)

	// Body is the default style for conversation text.
	Body = lipgloss.NewStyle().
		Foreground(White)

	// UserMessage styles messages sent by the user.
	UserMessage = lipgloss.NewStyle().
		Bold(true).
		Foreground(Secondary)

	// AssistantMessage styles messages from Nous itself.
	AssistantMessage = lipgloss.NewStyle().
		Foreground(Green)

	// Border is the base style used for viewport borders.
	Border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Muted)

	// Prompt styles the input area indicator (e.g. "> ").
	Prompt = lipgloss.NewStyle().
		Bold(true).
		Foreground(Accent)

	// Muted styles timestamps, metadata, and de-emphasised content.
	MutedStyle = lipgloss.NewStyle().
		Foreground(Muted)
)
