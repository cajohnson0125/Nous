// Package styles provides the single source of truth for all visual styling
// in Nous. Every color and base style used in the application is defined
// here using ANSI 16 colors only. No other package should use raw color
// values or ad-hoc styles — everything imports from this package.
package styles

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Standard ANSI 16 colors (indices 0–7).
//
// These map to the terminal's theme, so "Blue" is whatever the user's
// terminal defines for ANSI index 4, not a fixed hex shade.
var (
	Black   = lipgloss.Color("0")
	Red     = lipgloss.Color("1")
	Green   = lipgloss.Color("2")
	Yellow  = lipgloss.Color("3")
	Blue    = lipgloss.Color("4")
	Magenta = lipgloss.Color("5")
	Cyan    = lipgloss.Color("6")
	White   = lipgloss.Color("7")
)

// Bright ANSI 16 colors (indices 8–15).
var (
	BrightBlack  = lipgloss.Color("8")
	BrightRed    = lipgloss.Color("9")
	BrightGreen  = lipgloss.Color("10")
	BrightYellow = lipgloss.Color("11")
	BrightBlue   = lipgloss.Color("12")
	BrightMagenta = lipgloss.Color("13")
	BrightCyan   = lipgloss.Color("14")
	BrightWhite  = lipgloss.Color("15")
)

// Semantic color aliases.
//
// These provide meaningful names for the color roles used throughout
// Nous. Each maps to a standard ANSI 16 index — no custom values.
var (
	// Primary is the main accent color for headers and highlights.
	Primary = Blue

	// Secondary is used for supporting UI elements.
	Secondary = Cyan

	// Success indicates positive outcomes and confirmations.
	Success = Green

	// Warning draws attention to caution states.
	Warning = Yellow

	// Error highlights failures and problems.
	Error = Red

	// Muted is used for timestamps, metadata, and de-emphasised text.
	Muted = BrightBlack

	// Accent is for interactive elements like the prompt indicator.
	Accent = BrightCyan
)

// AllColors returns every exported color variable in this package,
// paired with its semantic name. This is used by tests to verify that
// all colors resolve to valid ANSI 16 indices.
func AllColors() map[string]color.Color {
	return map[string]color.Color{
		// Standard
		"Black":   Black,
		"Red":     Red,
		"Green":   Green,
		"Yellow":  Yellow,
		"Blue":    Blue,
		"Magenta": Magenta,
		"Cyan":    Cyan,
		"White":   White,
		// Bright
		"BrightBlack":   BrightBlack,
		"BrightRed":     BrightRed,
		"BrightGreen":   BrightGreen,
		"BrightYellow":  BrightYellow,
		"BrightBlue":    BrightBlue,
		"BrightMagenta": BrightMagenta,
		"BrightCyan":    BrightCyan,
		"BrightWhite":   BrightWhite,
		// Semantic
		"Primary":   Primary,
		"Secondary": Secondary,
		"Success":   Success,
		"Warning":   Warning,
		"Error":     Error,
		"Muted":     Muted,
		"Accent":    Accent,
	}
}
