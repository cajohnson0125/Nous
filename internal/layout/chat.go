// Package layout provides the chat TUI model for Nous.
// This is the primary user interface — a full-screen Bubbletea program
// with a scrollable conversation viewport and terminal-style input.
package layout

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/viewport"
	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"

	"github.com/cajohnson0125/Nous/internal/config"
	"github.com/cajohnson0125/Nous/internal/styles"
)

// Message represents a single message in the conversation.
type Message struct {
	Role    string // "user" or "assistant"
	Content string
}

// Model is the Bubbletea model for the chat TUI.
type Model struct {
	viewport viewport.Model
	input    string
	messages []Message
	cfg      *config.Config
	width    int
	height   int
	renderer *glamour.TermRenderer

	// headerHeight is the number of lines the header occupies.
	headerHeight int
	// inputHeight is the number of lines the input area occupies.
	inputHeight int
}

// NewModel creates a chat Model with the given config.
// The model is not fully initialized until it receives a WindowSizeMsg.
func NewModel(cfg *config.Config) Model {
	vp := viewport.New(
		viewport.WithWidth(80),
		viewport.WithHeight(24),
	)
	vp.Style = styles.Border.Copy().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.Muted)

	return Model{
		viewport:     vp,
		cfg:          cfg,
		headerHeight: 3, // title line + padding
		inputHeight:  1, // prompt line
	}
}

// Init returns the initial command.
// WindowSizeMsg is delivered automatically on program start.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleResize(msg)
		return m, nil

	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	// Delegate to viewport for scroll and mouse wheel handling.
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the full-screen chat TUI.
func (m Model) View() tea.View {
	header := m.renderHeader()

	viewportContent := m.viewport.View()

	inputLine := m.renderInput()

	// Compose the full layout: header + viewport + input.
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		viewportContent,
		inputLine,
	)

	// Calculate cursor position for the input line.
	// The cursor X is the visible width of the prompt prefix + input text.
	// The cursor Y is the total height of header + viewport + 1 (the input line).
	promptPrefix := styles.Prompt.Render("> ")
	cursorX := lipgloss.Width(promptPrefix) + lipgloss.Width(m.input)
	cursorY := m.headerHeight + lipgloss.Height(viewportContent) + 1

	cursor := m.buildCursor(cursorX, cursorY)

	v := tea.NewView(content)
	v.AltScreen = true
	v.WindowTitle = "Nous"
	v.Cursor = cursor

	return v
}

// handleResize updates dimensions and recalculates the viewport size.
func (m *Model) handleResize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height

	m.viewport.SetWidth(msg.Width)
	viewportHeight := msg.Height - m.headerHeight - m.inputHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	m.viewport.SetHeight(viewportHeight)

	// Recreate the Glamour renderer with updated width and theme.
	m.initRenderer()

	// Re-render viewport content after resize.
	m.updateViewportContent()
}

// handleKey processes key presses for typing, submission, quitting, and scrolling.
func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	s := msg.String()

	// Quit keys — always handled.
	if s == "ctrl+c" {
		return m, tea.Quit
	}
	if s == "q" && m.input == "" {
		return m, tea.Quit
	}

	// Scroll keys — only when input is empty.
	switch s {
	case "up", "k", "down", "j", "pgup", "pgdn":
		if m.input == "" {
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
		// When input is active, scroll keys are not consumed.
		return m, nil
	}

	// Submit.
	if s == "enter" {
		if m.input == "" {
			return m, nil
		}
		m.submitMessage()
		return m, nil
	}

	// Delete.
	if s == "backspace" {
		if len(m.input) > 0 {
			runes := []rune(m.input)
			m.input = string(runes[:len(runes)-1])
		}
		return m, nil
	}

	// Printable characters — append to input.
	if msg.Text != "" {
		m.input += msg.Text
	}

	return m, nil
}

// submitMessage adds the user message, generates the echo response, and clears input.
func (m *Model) submitMessage() {
	userMsg := Message{Role: "user", Content: m.input}
	m.messages = append(m.messages, userMsg)

	// Echo loop: echo back the user's message as an assistant response.
	echoMsg := Message{Role: "assistant", Content: m.input}
	m.messages = append(m.messages, echoMsg)

	m.input = ""
	m.updateViewportContent()
	m.viewport.GotoBottom()
}

// updateViewportContent rebuilds the viewport string from all messages.
func (m *Model) updateViewportContent() {
	content := FormatMessages(m.messages, m.renderer)
	m.viewport.SetContent(content)
}

// initRenderer creates a Glamour TermRenderer using the config theme
// and the current viewport width for word wrap.
func (m *Model) initRenderer() {
	theme := m.cfg.Theme
	if theme == "" {
		theme = "dark"
	}
	contentWidth := m.viewport.Width() - m.viewport.Style.GetHorizontalFrameSize()
	if contentWidth < 1 {
		contentWidth = 80
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath(theme),
		glamour.WithWordWrap(contentWidth),
	)
	if err != nil {
		// Fall back to no renderer on error.
		m.renderer = nil
		return
	}
	m.renderer = r
}

// renderHeader produces the centered "Nous" title bar.
func (m Model) renderHeader() string {
	title := styles.Title.Render("Nous")
	return lipgloss.Place(m.width, m.headerHeight, lipgloss.Center, lipgloss.Center, title)
}

// renderInput produces the prompt line with the current input text.
func (m Model) renderInput() string {
	prompt := styles.Prompt.Render("> ")
	text := styles.Body.Render(m.input)
	return prompt + text
}

// buildCursor creates a tea.Cursor based on config settings.
func (m Model) buildCursor(x, y int) *tea.Cursor {
	shape := cursorShape(m.cfg.Cursor)
	return &tea.Cursor{
		Position: tea.Position{X: x, Y: y},
		Shape:    shape,
		Blink:    m.cfg.Blink,
	}
}

// cursorShape maps a config cursor string to a tea.CursorShape.
func cursorShape(s string) tea.CursorShape {
	switch s {
	case "block":
		return tea.CursorBlock
	case "underline":
		return tea.CursorUnderline
	case "bar":
		return tea.CursorBar
	default:
		return tea.CursorBar
	}
}

// FormatMessages converts a slice of Messages into a styled string
// suitable for viewport content. The Glamour renderer is used for
// assistant messages when available; otherwise Lip Gloss styles are used.
func FormatMessages(msgs []Message, r *glamour.TermRenderer) string {
	if len(msgs) == 0 {
		return ""
	}

	var b strings.Builder
	for _, msg := range msgs {
		switch msg.Role {
		case "user":
			b.WriteString(styles.UserMessage.Render(fmt.Sprintf("You: %s", msg.Content)))
		case "assistant":
			content := msg.Content
			if r != nil {
				rendered, err := r.Render(content)
				if err == nil {
					content = rendered
				}
			}
			b.WriteString(styles.AssistantMessage.Render("Nous: "))
			b.WriteString(content)
		}
		b.WriteString("\n")
	}
	return b.String()
}
