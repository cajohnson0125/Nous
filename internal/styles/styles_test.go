package styles

import (
	"fmt"
	"image/color"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestAllColorsAreANSI16(t *testing.T) {
	colors := AllColors()

	for name, c := range colors {
		bc, ok := c.(ansi.BasicColor)
		if !ok {
			t.Errorf("color %q is not an ANSI 16 (BasicColor): got %T", name, c)
			continue
		}
		if bc > 15 {
			t.Errorf("color %q has ANSI index %d, want 0-15", name, bc)
		}
	}
}

func TestSemanticAliasesMapToStandardColors(t *testing.T) {
	tests := []struct {
		name   string
		color  color.Color
		expect string
	}{
		{"Primary", Primary, "Blue"},
		{"Secondary", Secondary, "Cyan"},
		{"Success", Success, "Green"},
		{"Warning", Warning, "Yellow"},
		{"Error", Error, "Red"},
		{"Muted", Muted, "BrightBlack"},
		{"Accent", Accent, "BrightCyan"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc, ok := tt.color.(ansi.BasicColor)
			if !ok {
				t.Fatalf("%s is not an ANSI BasicColor", tt.name)
			}

			var expected color.Color
			switch tt.expect {
			case "Blue":
				expected = Blue
			case "Cyan":
				expected = Cyan
			case "Green":
				expected = Green
			case "Yellow":
				expected = Yellow
			case "Red":
				expected = Red
			case "BrightBlack":
				expected = BrightBlack
			case "BrightCyan":
				expected = BrightCyan
			default:
				t.Fatalf("unhandled expected color: %s", tt.expect)
			}

			expectedBC, ok := expected.(ansi.BasicColor)
			if !ok {
				t.Fatalf("expected color %s is not an ANSI BasicColor", tt.expect)
			}
			if bc != expectedBC {
				t.Errorf("%s = ANSI %d (%s), want ANSI %d (%s)", tt.name, bc, colorName(bc), expectedBC, tt.expect)
			}
		})
	}
}

func TestStylesRenderWithoutPanic(t *testing.T) {
	tests := []struct {
		name  string
		style lipgloss.Style
		input string
	}{
		{"Title", Title, "Nous"},
		{"Body", Body, "Hello, world."},
		{"UserMessage", UserMessage, "User says hello"},
		{"AssistantMessage", AssistantMessage, "Nous responds"},
		{"Border", Border, "bordered content"},
		{"Prompt", Prompt, "> "},
		{"MutedStyle", MutedStyle, "2024-01-01 00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.style.Render(tt.input)
			if out == "" && tt.input != "" {
				t.Errorf("%s.Render(%q) returned empty string", tt.name, tt.input)
			}
		})
	}
}

func TestNoHexOrANSI256(t *testing.T) {
	colors := AllColors()
	for name, c := range colors {
		switch v := c.(type) {
		case ansi.BasicColor:
			if v > 15 {
				t.Errorf("%s is BasicColor but out of range: %d", name, v)
			}
		case ansi.IndexedColor:
			t.Errorf("%s is ANSI 256 (index %d), not ANSI 16", name, v)
		default:
			t.Errorf("%s is %T, expected ansi.BasicColor", name, v)
		}
	}
}

func colorName(c ansi.BasicColor) string {
	names := map[ansi.BasicColor]string{
		0:  "Black",
		1:  "Red",
		2:  "Green",
		3:  "Yellow",
		4:  "Blue",
		5:  "Magenta",
		6:  "Cyan",
		7:  "White",
		8:  "BrightBlack",
		9:  "BrightRed",
		10: "BrightGreen",
		11: "BrightYellow",
		12: "BrightBlue",
		13: "BrightMagenta",
		14: "BrightCyan",
		15: "BrightWhite",
	}
	if n, ok := names[c]; ok {
		return n
	}
	return fmt.Sprintf("Unknown(%d)", c)
}
