# Styling & Layout: Lip Gloss & Bubbles

> Research document for M001 -- Charmbracelet UI Stack Research
> Date: 2026-06-06
> Status: Phase 3

## Summary

This document covers **Lip Gloss** (charmbracelet/lipgloss) for declarative terminal styling and layout, and **Bubbles** (charmbracelet/bubbles) for pre-built Bubbletea components.

---

## 1. Lip Gloss (charmbracelet/lipgloss)

**Repository:** https://github.com/charmbracelet/lipgloss
**GoDoc:** https://pkg.go.dev/github.com/charmbracelet/lipgloss
**Latest Version:** `v2.0.3` (released 2026-04-13)
**License:** MIT
**Stars:** 11.4k | **Importers:** 9,630

> **Note:** Lip Gloss v2 uses import path `charm.land/lipgloss/v2`. The pkg.go.dev page for v1 shows `v1.1.0`; use the v2 path.

### 1.1 Role in Nous

Lip Gloss is the styling and layout engine. Every visual element in Nous's TUI -- borders, colors, padding, margins, alignment, tables, lists, trees -- is rendered through Lip Gloss styles. It is used inside Bubbletea `View()` methods to produce styled strings.

### 1.2 Style API

Styles are created declaratively using a fluent builder pattern:

```go
style := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    PaddingTop(2).
    PaddingLeft(4).
    Width(22)

fmt.Println(style.Render("Hello, kitty"))
```

**Key principle:** Styles are immutable value types. Every method call returns a new `Style`. Assignment creates a true copy.

### 1.3 Style Methods Reference

#### Text Formatting
| Method | Effect |
|--------|--------|
| `Bold(bool)` | Bold text |
| `Italic(bool)` | Italic text |
| `Faint(bool)` | Dimmed text |
| `Blink(bool)` | Blinking text |
| `Underline(bool)` | Underlined text |
| `Strikethrough(bool)` | Strikethrough text |
| `Reverse(bool)` | Swap fg/bg colors |
| `Transform(fn)` | Apply function (e.g., `strings.ToUpper`) |
| `Inline(bool)` | Single-line render, no margins/padding/borders |

#### Colors
| Method | Type | Purpose |
|--------|------|---------|
| `Foreground(c)` | `TerminalColor` | Text color |
| `Background(c)` | `TerminalColor` | Background color |

Color types:
- `lipgloss.Color("#FF0000")` -- True Color (24-bit hex)
- `lipgloss.Color("86")` -- ANSI 256 (8-bit)
- `lipgloss.Color("5")` -- ANSI 16 (4-bit)
- `lipgloss.AdaptiveColor{Light: "236", Dark: "248"}` -- Light/dark adaptive
- `lipgloss.CompleteColor{TrueColor: "#0000FF", ANSI256: "86", ANSI: "5"}` -- Profile-specific
- `lipgloss.CompleteAdaptiveColor{Light: ..., Dark: ...}` -- Combined
- `lipgloss.NoColor{}` -- No color (transparent)

#### Layout
| Method | Effect |
|--------|--------|
| `Width(int)` | Set block width (text wraps) |
| `Height(int)` | Set block height |
| `MaxWidth(int)` | Limit maximum width |
| `MaxHeight(int)` | Limit maximum height |
| `Padding(...int)` | Inner space (CSS shorthand: 1-4 values) |
| `Margin(...int)` | Outer space (CSS shorthand: 1-4 values) |
| `Align(...Position)` | Text alignment (Left, Center, Right) |
| `AlignHorizontal(Position)` | Horizontal alignment |
| `AlignVertical(Position)` | Vertical alignment |

#### Borders
| Method | Effect |
|--------|--------|
| `Border(b, sides...)` | Set border style and sides |
| `BorderStyle(Border)` | Border character set |
| `BorderForeground(c)` | Border text color |
| `BorderBackground(c)` | Border bg color |
| `BorderTop/Bottom/Left/Right(bool)` | Per-side toggle |

Border styles: `NormalBorder()`, `RoundedBorder()`, `BlockBorder()`, `ThickBorder()`, `DoubleBorder()`, `HiddenBorder()`, `InnerHalfBlockBorder()`, `OuterHalfBlockBorder()`, `ASCIIBorder()`, `MarkdownBorder()`.

#### Rendering
| Method | Effect |
|--------|--------|
| `Render(strs ...string) string` | Apply style to string |
| `String() string` | Stringer interface (needs `SetString`) |
| `SetString(strs ...string) Style` | Set default content |
| `Inherit(i Style) Style` | Copy unset rules from another style |
| `Renderer(r *Renderer) Style` | Bind to specific renderer |

### 1.4 Layout Utilities

#### Joining
```go
lipgloss.JoinHorizontal(lipgloss.Bottom, blockA, blockB, blockC)
lipgloss.JoinVertical(lipgloss.Center, blockA, blockB)
```

Position values: `Top` (0.0), `Center` (0.5), `Bottom` (1.0), `Left` (0.0), `Right` (1.0), or any float 0.0-1.0.

#### Measuring
```go
w := lipgloss.Width(block)
h := lipgloss.Height(block)
w, h := lipgloss.Size(block)
```

#### Placing
```go
lipgloss.PlaceHorizontal(80, lipgloss.Center, text)
lipgloss.PlaceVertical(30, lipgloss.Bottom, text)
lipgloss.Place(30, 80, lipgloss.Right, lipgloss.Bottom, text)
```

### 1.5 Sub-packages

#### `lipgloss/table` -- Table rendering
```go
t := table.New().
    Border(lipgloss.NormalBorder()).
    BorderStyle(lipgloss.NewStyle().Foreground(purple)).
    StyleFunc(func(row, col int) lipgloss.Style { ... }).
    Headers("LANGUAGE", "FORMAL", "INFORMAL").
    Rows(rows...)
fmt.Println(t)
```

#### `lipgloss/list` -- List rendering
```go
l := list.New("A", "B", "C").
    Enumerator(list.Roman).
    EnumeratorStyle(enumeratorStyle).
    ItemStyle(itemStyle)
```

#### `lipgloss/tree` -- Tree rendering
```go
t := tree.Root(".").
    Child("A", "B",
        tree.New().Root("Sub").Child("X", "Y"),
    )
```

### 1.6 Renderer

Custom renderers allow per-output color detection (critical for SSH):

```go
renderer := lipgloss.NewRenderer(sess) // sess is io.Writer
style := renderer.NewStyle().Background(lipgloss.AdaptiveColor{Light: "63", Dark: "228"})
```

### 1.7 Rendering Behavior

- Colors auto-downsample based on terminal capabilities
- Non-TTY output: all ANSI stripped
- Width/height account for East Asian wide characters
- Tabs converted to 4 spaces by default (configurable via `TabWidth`)
- `NO_COLOR` and `CLICOLOR_FORCE` environment variables respected

### 1.8 Minimal Working Example

```go
package main

import (
    "fmt"
    "strings"

    "charm.land/lipgloss/v2"
)

func main() {
    titleStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FAFAFA")).
        Background(lipgloss.Color("#7D56F4")).
        Padding(0, 2)

    boxStyle := lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#7D56F4")).
        Padding(1, 2).
        Width(40)

    title := titleStyle.Render("Nous Status")
    content := boxStyle.Render(
        "State: healthy\nNodes: 3/3\nVersion: 0.1.0",
    )

    fmt.Println(lipgloss.JoinVertical(lipgloss.Center, title, content))
}
```

### 1.9 Stability Assessment

- **Stable v2 release.** v2.0.3 released April 2026.
- **9,630 importers.** Industry-standard Go terminal styling library.
- **Mature API.** Expressive, well-documented, CSS-like paradigm.
- **Risk:** Very low. Lip Gloss is the foundation of the Charm ecosystem.

---

## 2. Bubbles (charmbracelet/bubbles)

**Repository:** https://github.com/charmbracelet/bubbles
**GoDoc:** https://pkg.go.dev/github.com/charmbracelet/bubbles
**Latest Version:** `v2.1.0` (released 2026-03-26)
**License:** MIT
**Stars:** 8.5k | **Forks:** 420

> **Note:** Bubbles v2 uses import path `charm.land/bubbles/v2`.

### 2.1 Role in Nous

Bubbles provides pre-built Bubbletea components that Nous uses for interactive UI elements: text inputs, tables, lists, spinners, progress bars, viewports, help views, and more. These are the building blocks for Nous's interactive screens.

### 2.2 Component Catalog

| Component | Sub-package | Purpose | Nous Usage |
|-----------|------------|---------|------------|
| **Spinner** | `spinner` | Loading/progress animation | Long-running operations, state transitions |
| **Text Input** | `textinput` | Single-line text field | Command input, search, filter |
| **Text Area** | `textarea` | Multi-line text editor | Configuration editing, notes |
| **Table** | `table` | Columnar data display | State listings, node status, diffs |
| **Progress** | `progress` | Progress bar | Apply progress, download status |
| **Paginator** | `paginator` | Pagination UI | Long lists, log pagination |
| **Viewport** | `viewport` | Scrollable content area | Log viewing, diff viewing, expert output |
| **List** | `list` | Browsable item list with filtering | Resource selection, command menus |
| **File Picker** | `filepicker` | File system browser | Config file selection |
| **Timer** | `timer` | Countdown timer | Approval gate timeouts |
| **Stopwatch** | `stopwatch` | Count-up timer | Operation timing |
| **Help** | `help` | Keybinding help view | All interactive screens |
| **Key** | `key` | Keybinding management | All interactive screens |

### 2.3 Component Integration Pattern

Every Bubble is a Bubbletea sub-model. It implements `Init()`, `Update()`, and `View()`. Integration requires embedding the Bubble in your model and delegating messages.

```go
type model struct {
    spinner spinner.Model
    table   table.Model
    viewport viewport.Model
    inputs  []textinput.Model
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd

    // Delegate to sub-models
    m.spinner, cmd = m.spinner.Update(msg)
    cmds = append(cmds, cmd)

    m.table, cmd = m.table.Update(msg)
    cmds = append(cmds, cmd)

    m.viewport, cmd = m.viewport.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}
```

### 2.4 Key Components Deep Dive

#### Spinner
```go
s := spinner.New()
s.Spinner = spinner.Dot
s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
```

#### Text Input
```go
ti := textinput.New()
ti.Placeholder = "Enter state name..."
ti.Focus()
ti.CharLimit = 50
ti.Width = 40
```

#### Table
```go
columns := []table.Column{
    {Title: "Name", Width: 20},
    {Title: "Status", Width: 10},
    {Title: "Version", Width: 10},
}
rows := []table.Row{{"node-1", "healthy", "0.1.0"}}
t := table.New(
    table.WithColumns(columns),
    table.WithRows(rows),
    table.WithFocused(true),
    table.WithHeight(7),
)
```

#### Viewport
```go
vp := viewport.New(width, height)
vp.SetContent(largeString)
vp.Style = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
```

#### List
```go
items := []list.Item{
    item{title: "Apply State", desc: "Apply a state configuration"},
    item{title: "Show Status", desc: "Show current state status"},
}
l := list.New(items, list.NewDefaultDelegate())
l.Title = "Nous Commands"
```

#### Help
```go
h := help.New()
h.Styles.ShortKey.Foreground(lipgloss.Color("63"))
// In View(): h.View(m.keys)
```

#### Key Bindings
```go
type keyMap struct {
    Up    key.Binding
    Down  key.Binding
    Enter key.Binding
    Quit  key.Binding
}

var keys = keyMap{
    Up: key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/↑", "move up")),
    Down: key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/↓", "move down")),
    Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
    Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}
```

### 2.5 Composition Patterns

**Pattern 1: Split pane with viewport + list**
```go
func (m model) View() tea.View {
    leftPane := m.list.View()
    rightPane := m.viewport.View()
    return *tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane))
}
```

**Pattern 2: Status bar + main content**
```go
func (m model) View() tea.View {
    statusBar := m.statusStyle.Render(fmt.Sprintf("Nodes: %d | Version: %s", m.nodes, m.version))
    main := m.viewport.View()
    return *tea.NewView(lipgloss.JoinVertical(lipgloss.Left, main, statusBar))
}
```

### 2.6 Stability Assessment

- **Stable v2 release.** v2.1.0 released March 2026.
- **8.5k stars**, used in production by Glow and many applications.
- **Active development.** v2.1.0 added dynamic height for textareas.
- **Bubbletea v2 compatible.** Import path `charm.land/bubbles/v2`.
- **Risk:** Low. Well-maintained, used by the Charm team's own products.

---

## 3. Pinned Versions

| Package | Import Path | Version | Date |
|---------|------------|---------|------|
| Lip Gloss | `charm.land/lipgloss/v2` | `v2.0.3` | 2026-04-13 |
| Bubbles | `charm.land/bubbles/v2` | `v2.1.0` | 2026-03-26 |

### go.mod directive

```
require (
    charm.land/lipgloss/v2 v2.0.3
    charm.land/bubbles/v2 v2.1.0
)
```

---

## 4. Gaps & Considerations

1. **Bubbles v2 compatibility** -- Verify all Bubbles components are v2-compatible at implementation time. The list component had significant v2 API changes.
2. **Table scrolling** -- Bubbles table supports vertical scrolling but not horizontal. For wide data, use viewport wrapping or implement custom horizontal scroll.
3. **Viewport + Lip Gloss** -- Use Lip Gloss styles on the viewport border, but the content area is raw strings. For rich content inside viewports, use Glamour (Phase 4) or manual Lip Gloss styling.
4. **Custom delegates** -- List items need a custom `list.Item` implementation and delegate for styled rendering. Plan for Nous-specific item types.
5. **Key binding conflicts** -- When nesting components (viewport inside list), key bindings can conflict. Use `tea.WithFilter()` or component-level key routing.

---

## References

- Lip Gloss GitHub: https://github.com/charmbracelet/lipgloss
- Lip Gloss GoDoc: https://pkg.go.dev/github.com/charmbracelet/lipgloss
- Bubbles GitHub: https://github.com/charmbracelet/bubbles
- Bubbles GoDoc: https://pkg.go.dev/github.com/charmbracelet/bubbles
