# User Interaction: Huh & Glamour

> Research document for M001 -- Charmbracelet UI Stack Research
> Date: 2026-06-06
> Status: Phase 4

## Summary

This document covers **Huh** (charmbracelet/huh) for interactive forms and prompts, and **Glamour** (charmbracelet/glamour) for stylesheet-based Markdown rendering.

---

## 1. Huh (charmbracelet/huh)

**Repository:** https://github.com/charmbracelet/huh
**GoDoc:** https://pkg.go.dev/github.com/charmbracelet/huh
**Latest Version:** `v2.0.3` (released 2026-03-10)
**License:** MIT
**Stars:** 6.9k | **Forks:** 247

> **Note:** Huh v2 uses import path `charm.land/huh/v2`.

### 1.1 Role in Nous

Huh provides interactive forms, prompts, and surveys. In Nous, Huh handles human-in-the-loop approval gates, configuration wizards, and confirmation dialogs. It is a Bubbletea program under the hood.

### 1.2 Core API

#### Form

A `Form` is a collection of `Group`s, executed sequentially.

```go
form := huh.NewForm(
    huh.NewGroup(
        huh.NewInput().Title("State name").Value(&name),
        huh.NewSelect[string]().Title("Environment").
            Options(huh.NewOptions("dev", "staging", "prod")...).
            Value(&env),
    ),
    huh.NewGroup(
        huh.NewConfirm().Title("Apply this configuration?").Value(&confirmed),
    ),
)
err := form.Run()
```

#### Field Types

| Field | Constructor | Purpose |
|-------|------------|---------|
| **Input** | `huh.NewInput()` | Single-line text input |
| **Text** | `huh.NewText()` | Multi-line text area |
| **Select** | `huh.NewSelect[T]()` | Single-select from list |
| **MultiSelect** | `huh.NewMultiSelect[T]()` | Multi-select from list |
| **Confirm** | `huh.NewConfirm()` | Yes/No confirmation |
| **FilePicker** | `huh.NewFilePicker()` | File system browser |

#### Field Methods (common)

| Method | Purpose |
|--------|---------|
| `.Title(string)` | Set field title |
| `.Description(string)` | Set description |
| `.Value(&var)` | Bind to variable |
| `.Key(string)` | Unique key for access |
| `.Default(val)` | Default value |
| `.Validate(fn)` | Validation function |
| `.Inline(true)` | Inline layout |
| `.WithTheme(theme)` | Custom theme |
| `.WithWidth(int)` | Set width |
| `.WithHeight(int)` | Set height |
| `.WithButtonTag(tag)` | Custom button text (Confirm) |
| `.WithHelp(help)` | Custom help text |
| `.CharLimit(int)` | Max characters (Input) |

#### Select/MultiSelect Options

```go
options := huh.NewOptions("dev", "staging", "prod")
// Or with values:
options := []huh.Option[string]{
    huh.NewOption("Development", "dev"),
    huh.NewOption("Staging", "staging"),
    huh.NewOption("Production", "prod"),
}
```

### 1.3 Themes

```go
theme := huh.ThemeCharm()
// Or: huh.ThemeBase(), huh.ThemeDracula(), huh.ThemeCatppuccin()
form.WithTheme(theme)
```

Custom themes use Lip Gloss styles:

```go
theme := &huh.Theme{
    Form: huh.FormStyles{
        Selector: lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")),
    },
}
```

### 1.4 Bubbletea Integration

Huh forms can run standalone (`form.Run()`) or embedded in Bubbletea programs:

```go
// Standalone (blocks until complete)
err := form.Run()

// Embedded in Bubbletea
field := huh.NewInput().Title("Name")
field.Init()
// In Update: field.Update(msg)
// In View: field.View()
```

### 1.5 Accessibility

Huh supports `WithAccessible(true)` for screen-reader-friendly mode:

```go
form := huh.NewForm(group...).WithAccessible(true)
```

### 1.6 Minimal Working Example

```go
package main

import (
    "fmt"
    "log"

    "charm.land/huh/v2"
)

func main() {
    var name string
    var env string
    var confirmed bool

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("State name").
                Placeholder("my-state").
                Value(&name),
            huh.NewSelect[string]().
                Title("Target environment").
                Options(
                    huh.NewOption("Development", "dev"),
                    huh.NewOption("Staging", "staging"),
                    huh.NewOption("Production", "prod"),
                ).
                Value(&env),
        ),
        huh.NewGroup(
            huh.NewConfirm().
                Title(fmt.Sprintf("Apply %q to %s?", name, env)).
                Value(&confirmed),
        ),
    )

    if err := form.Run(); err != nil {
        log.Fatal(err)
    }

    if confirmed {
        fmt.Printf("Applying %s to %s\n", name, env)
    }
}
```

### 1.7 Stability Assessment

- **Stable v2 release.** v2.0.3 released March 2026.
- **6.9k stars**, active development.
- **Well-tested** in production (Glow uses it).
- **Risk:** Low. Clean API, stable semver.

---

## 2. Glamour (charmbracelet/glamour)

**Repository:** https://github.com/charmbracelet/glamour
**GoDoc:** https://pkg.go.dev/github.com/charmbracelet/glamour
**Latest Version:** `v2.0.0` (released 2026)
**License:** MIT
**Stars:** 2.9k

> **Note:** Glamour v2 uses import path `charm.land/glamour/v2`. This is a vanity domain redirecting to the GitHub repository.

### 2.1 Role in Nous

Glamour renders Markdown content as styled terminal output. In Nous, it handles expert output rendering (reports, documentation, status details), long-form content in viewports, and any Markdown that needs rich terminal presentation.

### 2.2 Core API

#### Simple Rendering

```go
import "charm.land/glamour/v2"

out, err := glamour.Render("# Hello\n\nWorld", "dark")
fmt.Print(out)
```

#### TermRenderer (recommended for multiple renders)

```go
r, err := glamour.NewTermRenderer(
    glamour.WithStylePath("dark"),
    glamour.WithWordWrap(80),
)

out1, _ := r.Render("# Title 1\nContent")
out2, _ := r.Render("# Title 2\nMore content")
```

#### Render Options

| Option | Purpose |
|--------|---------|
| `WithStylePath(path)` | Style to use: "dark", "light", "pink", "dracula", "tokyo-night", or custom |
| `WithWordWrap(width)` | Word wrap width (0 = disabled) |
| `WithPreservedNewLines()` | Preserve newlines in output |
| `WithCustomRenderer(r)` | Custom Goldmark renderer |
| `WithStyles(s)` | Custom style definitions |
| `WithBaseURL(url)` | Base URL for relative links |
| `WithEmoji()` | Enable emoji rendering |

### 2.3 Built-in Styles

- `"dark"` -- Default dark theme
- `"light"` -- Light background theme
- `"pink"` -- Pink theme
- `"dracula"` -- Dracula color scheme
- `"tokyo-night"` -- Tokyo Night theme
- `"notty"` -- No terminal formatting (plain text)

### 2.4 v2 Changes from v1

| Change | Impact |
|--------|--------|
| Import path: `charm.land/glamour/v2` | Module path changed |
| Lip Gloss v2 integration | Better color handling |
| `WithAutoStyle()` removed | Default is now "dark" |
| `WithColorProfile()` removed | Use `lipgloss.Print()` for downsampling |
| `Overlined` style removed | Not widely supported |
| Hyperlink support (OSC 8) | Clickable links in supported terminals |
| Better CJK/emoji wrapping | Uses `lipgloss.Wrap` |
| Email autolinks hide `mailto:` | Cleaner rendering |

### 2.5 Viewport Integration

Glamour output is a string. For scrollable display in a TUI, use with Bubbles viewport:

```go
// Render markdown
renderer, _ := glamour.NewTermRenderer(glamour.WithWordWrap(width))
content, _ := renderer.Render(markdownString)

// Display in viewport
vp := viewport.New(width, height)
vp.SetContent(content)
```

### 2.6 Custom Styles

Styles are defined as JSON or YAML:

```go
// Custom style JSON
customStyle := `
{
  "document": {
    "margin": 2
  },
  "h1": {
    "margin_top": 2,
    "margin_bottom": 1,
    "bold": true,
    "color": "#7D56F4"
  }
}
`
r, _ := glamour.NewTermRenderer(
    glamour.WithStylesFromJSON(strings.NewReader(customStyle)),
)
```

### 2.7 Minimal Working Example

```go
package main

import (
    "fmt"

    "charm.land/glamour/v2"
)

func main() {
    md := `# Nous Status Report

## Overview
All systems are **operational**.

### Node Status
| Node    | Status  | Uptime  |
|---------|---------|---------|
| node-1  | healthy | 99.9%   |
| node-2  | healthy | 99.8%   |
| node-3  | healthy | 99.7%   |

### Recent Events
- State applied successfully
- Configuration drift detected and corrected
- [View full report](https://nous.example.com/report)

> **Note:** Next maintenance window is Saturday 02:00 UTC.
`

    out, err := glamour.Render(md, "dark")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Print(out)
}
```

### 2.8 Stability Assessment

- **Stable v2 release.** v2.0.0 is a major but well-documented upgrade.
- **Breaking changes from v1** -- import path, removed options. Upgrade guide available.
- **Used in production** by Glow and many CLI tools.
- **Risk:** Low. Glamour is focused and mature. The v2 API is cleaner.

---

## 3. Pinned Versions

| Package | Import Path | Version | Date |
|---------|------------|---------|------|
| Huh | `charm.land/huh/v2` | `v2.0.3` | 2026-03-10 |
| Glamour | `charm.land/glamour/v2` | `v2.0.0` | 2026 |

### go.mod directive

```
require (
    charm.land/huh/v2 v2.0.3
    charm.land/glamour/v2 v2.0.0
)
```

---

## 4. Gaps & Considerations

1. **Glamour v2 import path** -- Uses `charm.land/glamour/v2` vanity domain. Verify go.sum resolves correctly.
2. **Glamour color downsampling** -- Glamour v2 does NOT auto-downsample colors. Use `lipgloss.Print(out)` to downsample for the terminal. This is a change from v1.
3. **Huh + Bubbletea** -- For complex flows (e.g., approval gate with countdown timer), embed Huh fields inside a custom Bubbletea model rather than using `form.Run()`.
4. **Viewport width** -- Glamour needs to know the viewport width for word wrap. Pass `glamour.WithWordWrap(vp.Width - frameSize)` dynamically.
5. **Accessibility** -- Always offer `WithAccessible(true)` as a fallback for screen readers and non-interactive environments.

---

## References

- Huh GitHub: https://github.com/charmbracelet/huh
- Huh GoDoc: https://pkg.go.dev/github.com/charmbracelet/huh
- Glamour GitHub: https://github.com/charmbracelet/glamour
- Glamour GoDoc: https://pkg.go.dev/github.com/charmbracelet/glamour
