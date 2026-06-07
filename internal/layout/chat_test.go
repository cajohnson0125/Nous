package layout_test

import (
	"bytes"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"

	"github.com/cajohnson0125/Nous/internal/config"
	"github.com/cajohnson0125/Nous/internal/layout"
)

func defaultModel() layout.Model {
	return layout.NewModel(config.Default())
}

// updateModel calls Update and type-asserts the result back to layout.Model.
// Panics if the type assertion fails.
func updateModel(m layout.Model, msg tea.Msg) (layout.Model, tea.Cmd) {
	model, cmd := m.Update(msg)
	return model.(layout.Model), cmd
}

// --- FormatMessages tests ---

func TestFormatMessagesEmpty(t *testing.T) {
	result := layout.FormatMessages(nil, nil)
	if result != "" {
		t.Errorf("expected empty string for nil messages, got %q", result)
	}

	result = layout.FormatMessages([]layout.Message{}, nil)
	if result != "" {
		t.Errorf("expected empty string for empty messages, got %q", result)
	}
}

func TestFormatMessagesUserOnly(t *testing.T) {
	msgs := []layout.Message{
		{Role: "user", Content: "Hello"},
	}
	result := layout.FormatMessages(msgs, nil)

	if !strings.Contains(result, "You: Hello") {
		t.Errorf("expected user message in output, got %q", result)
	}
	if !strings.HasSuffix(result, "\n") {
		t.Errorf("expected trailing newline, got %q", result)
	}
}

func TestFormatMessagesUserAndAssistant(t *testing.T) {
	msgs := []layout.Message{
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there"},
	}
	result := layout.FormatMessages(msgs, nil)

	if !strings.Contains(result, "You: Hello") {
		t.Errorf("expected user message, got %q", result)
	}
	if !strings.Contains(result, "Nous:") {
		t.Errorf("expected assistant prefix, got %q", result)
	}
	if !strings.Contains(result, "Hi there") {
		t.Errorf("expected assistant content, got %q", result)
	}
}

func TestFormatMessagesWithGlamourRenderer(t *testing.T) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("create glamour renderer: %v", err)
	}

	msgs := []layout.Message{
		{Role: "assistant", Content: "Hello **world**"},
	}
	result := layout.FormatMessages(msgs, r)

	if !strings.Contains(result, "Nous:") {
		t.Errorf("expected assistant prefix, got %q", result)
	}
	// Glamour renders **bold** with ANSI escape sequences, not raw markdown.
	if strings.Contains(result, "**world**") {
		t.Errorf("expected markdown to be rendered, but raw markdown found in %q", result)
	}
}

// --- cursorShape tests ---

func TestCursorShape(t *testing.T) {
	tests := []struct {
		name  string
		input string
		wants tea.CursorShape
	}{
		{"block", "block", tea.CursorBlock},
		{"underline", "underline", tea.CursorUnderline},
		{"bar", "bar", tea.CursorBar},
		{"unknown defaults to bar", "unknown", tea.CursorBar},
		{"empty defaults to bar", "", tea.CursorBar},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Cursor: tt.input}
			m := layout.NewModel(cfg)
			// Trigger resize so the view renders properly.
			m, _ = updateModel(m, tea.WindowSizeMsg{Width: 80, Height: 24})
			v := m.View()
			if v.Cursor == nil {
				t.Fatal("expected cursor to be non-nil")
			}
			if v.Cursor.Shape != tt.wants {
				t.Errorf("cursor shape = %v, want %v", v.Cursor.Shape, tt.wants)
			}
		})
	}
}

// --- Key handling tests ---

func TestKeyTyping(t *testing.T) {
	m := defaultModel()

	// Type "hello"
	for _, ch := range "hello" {
		msg := tea.KeyPressMsg{Text: string(ch), Code: rune(ch)}
		m, _ = updateModel(m, msg)
	}

	v := m.View()
	content := v.Content
	if !strings.Contains(content, "hello") {
		t.Errorf("expected 'hello' in view, got %q", content)
	}
}

func TestKeyBackspace(t *testing.T) {
	m := defaultModel()

	// Type "hi"
	for _, ch := range "hi" {
		msg := tea.KeyPressMsg{Text: string(ch), Code: rune(ch)}
		m, _ = updateModel(m, msg)
	}

	// Press backspace
	msg := tea.KeyPressMsg{Code: tea.KeyBackspace}
	m, _ = updateModel(m, msg)

	v := m.View()
	content := v.Content
	// After backspace, "hello" should not appear — only "h" remains.
	if strings.Contains(content, "hello") {
		t.Errorf("expected 'h' after backspace, not 'hello', got %q", content)
	}
}

func TestKeyBackspaceEmpty(t *testing.T) {
	m := defaultModel()

	// Press backspace on empty input — should not panic.
	msg := tea.KeyPressMsg{Code: tea.KeyBackspace}
	m, _ = updateModel(m, msg)

	v := m.View()
	if v.Content == "" {
		t.Error("expected non-empty view even with empty input")
	}
}

func TestKeyEnterSubmit(t *testing.T) {
	m := defaultModel()

	// Type "hello"
	for _, ch := range "hello" {
		msg := tea.KeyPressMsg{Text: string(ch), Code: rune(ch)}
		m, _ = updateModel(m, msg)
	}

	// Press enter
	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	m, _ = updateModel(m, msg)

	// The viewport should now contain the user message and echo response.
	v := m.View()
	content := v.Content
	if !strings.Contains(content, "You: hello") {
		t.Errorf("expected 'You: hello' in viewport after submit, got %q", content)
	}
	if !strings.Contains(content, "Nous:") {
		t.Errorf("expected 'Nous:' echo in viewport after submit, got %q", content)
	}
}

func TestKeyEnterEmpty(t *testing.T) {
	m := defaultModel()

	// Press enter on empty input — should be a no-op.
	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	m, _ = updateModel(m, msg)

	v := m.View()
	content := v.Content
	if strings.Contains(content, "You:") {
		t.Error("expected no user message on empty enter")
	}
}

func TestKeyQuitCtrlC(t *testing.T) {
	m := defaultModel()

	// ctrl+c: Code is 'c' with ModCtrl modifier.
	msg := tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
	_, cmd := updateModel(m, msg)

	if cmd == nil {
		t.Error("expected non-nil command for ctrl+c (quit)")
	}
}

func TestKeyQuitQEmpty(t *testing.T) {
	m := defaultModel()

	msg := tea.KeyPressMsg{Text: "q", Code: 'q'}
	_, cmd := updateModel(m, msg)

	if cmd == nil {
		t.Error("expected non-nil command for 'q' with empty input (quit)")
	}
}

func TestKeyQWithInput(t *testing.T) {
	m := defaultModel()

	// Type "a" first so input is not empty
	msg := tea.KeyPressMsg{Text: "a", Code: 'a'}
	m, _ = updateModel(m, msg)

	// Press "q" — should NOT quit, should type 'q'
	msg = tea.KeyPressMsg{Text: "q", Code: 'q'}
	m, cmd := updateModel(m, msg)

	if cmd != nil {
		t.Error("'q' should produce no command when input is non-empty")
	}

	v := m.View()
	content := v.Content
	if !strings.Contains(content, "aq") {
		t.Errorf("expected 'aq' in view after typing 'a' then 'q', got %q", content)
	}
}

func TestEchoLoop(t *testing.T) {
	m := defaultModel()

	// Type and submit a message
	for _, ch := range "test message" {
		msg := tea.KeyPressMsg{Text: string(ch), Code: rune(ch)}
		m, _ = updateModel(m, msg)
	}

	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	m, _ = updateModel(m, msg)

	v := m.View()
	content := v.Content

	if !strings.Contains(content, "You: test message") {
		t.Errorf("expected user message after echo, got %q", content)
	}
	if !strings.Contains(content, "Nous:") {
		t.Errorf("expected assistant prefix after echo, got %q", content)
	}
	if !strings.Contains(content, "test message") {
		t.Errorf("expected echo content after submit, got %q", content)
	}
}

func TestMultipleSubmits(t *testing.T) {
	m := defaultModel()

	// Submit first message
	for _, ch := range "first" {
		msg := tea.KeyPressMsg{Text: string(ch), Code: rune(ch)}
		m, _ = updateModel(m, msg)
	}
	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	m, _ = updateModel(m, msg)

	// Submit second message
	for _, ch := range "second" {
		msg := tea.KeyPressMsg{Text: string(ch), Code: rune(ch)}
		m, _ = updateModel(m, msg)
	}
	msg = tea.KeyPressMsg{Code: tea.KeyEnter}
	m, _ = updateModel(m, msg)

	v := m.View()
	content := v.Content

	if !strings.Contains(content, "You: first") {
		t.Errorf("expected first user message, got %q", content)
	}
	if !strings.Contains(content, "You: second") {
		t.Errorf("expected second user message, got %q", content)
	}
}

// --- Quit via full program test ---

func TestProgramQuitsOnCtrlC(t *testing.T) {
	var buf bytes.Buffer
	var in bytes.Buffer

	cfg := config.Default()
	m := layout.NewModel(cfg)

	p := tea.NewProgram(m,
		tea.WithInput(&in),
		tea.WithOutput(&buf),
		tea.WithoutRenderer(),
	)

	// Inject ctrl+c via Send.
	go func() {
		p.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	}()

	_, err := p.Run()
	if err != nil {
		t.Errorf("expected clean exit, got error: %v", err)
	}
}

// --- Resize tests ---

func TestResize(t *testing.T) {
	m := defaultModel()

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	m, _ = updateModel(m, msg)

	v := m.View()
	content := v.Content
	if content == "" {
		t.Error("expected non-empty view after resize")
	}
	if !v.AltScreen {
		t.Error("expected AltScreen to be true")
	}
	if v.WindowTitle != "Nous" {
		t.Errorf("expected WindowTitle 'Nous', got %q", v.WindowTitle)
	}
}

func TestResizeSmall(t *testing.T) {
	m := defaultModel()

	// Very small window — should not crash.
	msg := tea.WindowSizeMsg{Width: 10, Height: 5}
	m, _ = updateModel(m, msg)

	v := m.View()
	if v.Content == "" {
		t.Error("expected non-empty view even with small window")
	}
}

// --- Blink config test ---

func TestBlinkFromConfig(t *testing.T) {
	cfg := &config.Config{
		Cursor: "bar",
		Blink:  true,
	}
	m := layout.NewModel(cfg)
	m, _ = updateModel(m, tea.WindowSizeMsg{Width: 80, Height: 24})

	v := m.View()
	if v.Cursor == nil {
		t.Fatal("expected cursor to be non-nil")
	}
	if !v.Cursor.Blink {
		t.Error("expected cursor blink to be true")
	}
}

func TestNoBlinkFromConfig(t *testing.T) {
	cfg := &config.Config{
		Cursor: "block",
		Blink:  false,
	}
	m := layout.NewModel(cfg)
	m, _ = updateModel(m, tea.WindowSizeMsg{Width: 80, Height: 24})

	v := m.View()
	if v.Cursor == nil {
		t.Fatal("expected cursor to be non-nil")
	}
	if v.Cursor.Blink {
		t.Error("expected cursor blink to be false")
	}
}
