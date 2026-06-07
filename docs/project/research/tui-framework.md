# TUI Framework: Bubbletea, Ultraviolet, Colorprofile

> Research document for M001 -- Charmbracelet UI Stack Research
> Date: 2026-06-06
> Status: Phase 2

## Summary

This document covers the three packages at the core of Nous's TUI layer: **Bubbletea** (charmbracelet/bubbletea) as the Elm-architecture TUI framework, **Ultraviolet** (charmbracelet/ultraviolet) as the low-level terminal primitives powering Bubbletea v2, and **Colorprofile** (charmbracelet/colorprofile) for color detection and downsampled output.

---

## 1. Bubbletea (charmbracelet/bubbletea)

**Repository:** https://github.com/charmbracelet/bubbletea
**GoDoc:** https://pkg.go.dev/github.com/charmbracelet/bubbletea
**Latest Version:** `v2.0.7` (released 2026-06-01)
**License:** MIT
**Stars:** 42.9k | **Forks:** 1.2k

> **Note:** Bubbletea v2 uses import path `charm.land/bubbletea/v2` as shown in the official README tutorial. The GitHub-hosted pkg.go.dev path shows v1 API. The v2 module path is `charm.land/bubbletea/v2`.

### 1.1 Role in Nous

Bubbletea is the TUI framework. Every interactive terminal screen in Nous -- dashboards, approval gates, expert output viewers, interactive workflows -- is a Bubbletea program. It provides the Elm Architecture (Model-Update-View) pattern that Nous's state engine maps onto naturally.

### 1.2 Architecture: The Elm Architecture

Bubbletea programs are built around three core functions on a `Model`:

```go
type Model interface {
    Init() Cmd
    Update(Msg) (Model, Cmd)
    View() View
}
```

- **Init()** -- Returns an initial command (or nil). Called once when the program starts.
- **Update(Msg)** -- Receives messages (key presses, mouse events, timer ticks, custom messages) and returns an updated model plus an optional command.
- **View()** -- Renders the current model state as a `tea.View`. Called after every Update.

**Data flow:**

```
Init() -> Cmd -> Msg -> Update() -> (Model, Cmd) -> View() -> render
                                        |
                                        +-> Cmd -> Msg -> Update() -> ...
```

### 1.3 Core Types

#### `Model` interface

```go
type Model interface {
    Init() Cmd
    Update(msg Msg) (Model, Cmd)
    View() View
}
```

#### `Msg` type

```go
type Msg interface{} // Any type. Messages trigger Update.
```

Built-in message types:
- `KeyPressMsg` -- Keyboard input (replaces v1's `KeyMsg`)
- `MouseMsg` -- Mouse events (click, scroll, motion)
- `WindowSizeMsg` -- Terminal resize: `{Width, Height int}`
- `FocusMsg` / `BlurMsg` -- Terminal focus gained/lost
- `QuitMsg` -- Program should quit
- `InterruptMsg` -- SIGINT received
- `SuspendMsg` -- SIGTSTP received

#### `Cmd` type

```go
type Cmd func() Msg // An IO operation that returns a message.
```

Built-in commands:

| Command | Purpose |
|---------|---------|
| `Batch(cmds ...Cmd) Cmd` | Run commands concurrently |
| `Sequence(cmds ...Cmd) Cmd` | Run commands sequentially, in order |
| `Tick(d, fn) Cmd` | Timer at interval (not clock-synced) |
| `Every(d, fn) Cmd` | Timer synced to system clock |
| `ExecProcess(cmd, fn) Cmd` | Run an exec.Cmd (blocking, e.g., editor) |
| `Exec(c, fn) Cmd` | Custom blocking IO (v0.21.0+) |
| `Printf(template, args...) Cmd` | Print above the program |
| `Println(args...) Cmd` | Print above the program |
| `SetWindowTitle(title) Cmd` | Set terminal title |
| `WindowSize() Cmd` | Query terminal size |
| `Quit` | Quit the program (special sentinel) |

#### Commands → View Fields (v2 change)

In v2, several v1 commands are now **declarative View fields** instead of commands:

| v1 Command | v2 View Field |
|------------|---------------|
| `tea.EnterAltScreen` | `view.AltScreen = true` |
| `tea.ExitAltScreen` | `view.AltScreen = false` |
| `tea.EnableMouseCellMotion` | `view.MouseMode = tea.MouseModeCellMotion` |
| `tea.EnableMouseAllMotion` | `view.MouseMode = tea.MouseModeAllMotion` |
| `tea.HideCursor` | `view.Cursor = nil` |
| `tea.ShowCursor` | `view.Cursor = &tea.Cursor{...}` or `tea.NewCursor(x, y)` |
| `tea.SetWindowTitle("...")` | `view.WindowTitle = "..."` |

```go
// v2 example: set features declaratively in View()
func (m model) View() tea.View {
    v := tea.NewView("Hello, world!")
    v.AltScreen = true
    v.MouseMode = tea.MouseModeCellMotion
    v.WindowTitle = "Nous"
    return *v
}
```

#### `Key` and `KeyPressMsg`

```go
type Key struct {
    Code  KeyCode    // KeyEnter, KeyRunes, KeyCtrlC, etc.
    Text  string     // Text content for KeyRunes
    Mod   Modifier   // Modifier keys (Alt, Shift, Ctrl, etc.)
}

type KeyPressMsg = Key  // v2 key message type

// KeyType constants:
KeyRunes, KeyUp, KeyDown, KeyLeft, KeyRight,
KeyEnter, KeyBackspace, KeyTab, KeyEsc, KeySpace, KeyDelete,
KeyHome, KeyEnd, KeyPgUp, KeyPgDown,
KeyF1...KeyF20,
KeyCtrlA...KeyCtrlZ, KeyCtrlC, etc.
KeyShiftUp/Down/Left/Right, KeyCtrlUp/Down/Left/Right,
KeyCtrlShiftUp/Down/Left/Right, KeyCtrlHome/End
```

### 1.4 Program and Options

```go
p := tea.NewProgram(model, opts...)
model, err := p.Run()
```

| ProgramOption | Purpose |
|---------------|---------|
| `WithContext(ctx)` | Context for cancellation |
| `WithInput(r io.Reader)` | Custom input (default: stdin) |
| `WithInputTTY()` | Open a new TTY for input |
| `WithOutput(w io.Writer)` | Custom output (default: stdout) |
| `WithEnvironment(env)` | Environment variables |
| `WithFPS(fps int)` | Custom max FPS (default 60, max 120) |
| `WithFilter(fn)` | Intercept messages before Update |
| `WithoutSignalHandler()` | Handle signals yourself |
| `WithoutSignals()` | Ignore OS signals (testing) |
| `WithoutRenderer()` | Non-TUI mode (no rendering) |
| `WithoutCatchPanics()` | Disable panic recovery |

> **v2 change:** Alt screen, mouse modes, focus reporting, bracketed paste, cursor, and window title are now **declarative View fields** set in `View()` rather than ProgramOptions or commands. See the View Fields table below.

#### View Fields (declarative, v2)

| Field | Type | What It Does |
|-------|------|-------------|
| `Content` | `string` | Set via `SetContent()` or `NewView()` |
| `AltScreen` | `bool` | Enter/exit alt screen buffer |
| `MouseMode` | `MouseMode` | None/CellMotion/AllMotion |
| `ReportFocus` | `bool` | Focus/blur events |
| `DisableBracketedPasteMode` | `bool` | Disable bracketed paste |
| `WindowTitle` | `string` | Terminal window title |
| `Cursor` | `*Cursor` | Position, shape, color, blink |
| `ForegroundColor` | `Color` | Terminal foreground |
| `BackgroundColor` | `Color` | Terminal background |

#### Program Methods

| Method | Purpose |
|--------|---------|
| `Run() (Model, error)` | Start and block until quit |
| `Send(msg Msg)` | Inject message from outside |
| `Quit()` | Quit from outside |
| `Kill()` | Force kill, skip final render |
| `ReleaseTerminal()` | Give terminal back (for shelling out) |
| `RestoreTerminal()` | Reclaim terminal after release |
| `Printf(template, args...)` | Print above program |
| `Println(args...)` | Print above program |
| `Wait()` | Block until program finishes |

### 1.5 Minimal Working Example

```go
package main

import (
    "fmt"
    "os"

    tea "charm.land/bubbletea/v2"
)

type model struct {
    cursor   int
    choices  []string
    selected map[int]struct{}
}

func initialModel() model {
    return model{
        choices:  []string{"Apply state", "Show status", "Inspect diff"},
        selected: make(map[int]struct{}),
    }
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        case "enter", "space":
            if _, ok := m.selected[m.cursor]; ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        }
    }
    return m, nil
}

func (m model) View() tea.View {
    s := "What would you like to do?\n\n"
    for i, choice := range m.choices {
        cursor := " "
        if m.cursor == i {
            cursor = ">"
        }
        checked := " "
        if _, ok := m.selected[i]; ok {
            checked = "x"
        }
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }
    s += "\nPress q to quit.\n"
    return *tea.NewView(s)
}

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
}
```

### 1.6 Key Changes from v1 to v2

- **KeyMsg renamed to KeyPressMsg** -- `KeyMsg` is still available as alias
- **View returns `tea.View`** -- a new wrapper type. Use `tea.NewView("content")` or `var v tea.View; v.SetContent("content")`.
- **Import path** is `charm.land/bubbletea/v2`
- **ProgramOptions → View fields** -- `WithAltScreen()`, `WithMouseCellMotion()`, `WithReportFocus()`, `WithoutBracketedPaste()` are now declarative View fields
- **Commands → View fields** -- `EnterAltScreen`, `ExitAltScreen`, `EnableMouseCellMotion`, `HideCursor`, `ShowCursor`, `SetWindowTitle` are now declarative View fields
- **Key fields changed** -- `msg.Type` → `msg.Code`, `msg.Runes` → `msg.Text` (string), `msg.Alt` → `msg.Mod.Contains(tea.ModAlt)`
- **p.Start() / p.StartReturningModel()** → `p.Run()`
- **Ultraviolet** replaces v1's internal terminal primitives
- **Scrolling commands deprecated** -- `ScrollDown`, `ScrollUp`, `SyncScrollArea` are deprecated

### 1.7 Stability Assessment

- **Stable v2 release.** v2.0.7 released June 1, 2026.
- **42.9k stars**, 18,000+ applications built with it.
- **Used in production** by Microsoft Azure, CockroachDB, AWS, NVIDIA, MinIO, Ubuntu.
- **Active development** -- 101 open issues, regular patch releases.
- **Breaking changes from v1 to v2** -- upgrade guide available.
- **Risk:** Very low. Industry-standard Go TUI framework.

---

## 2. Ultraviolet (charmbracelet/ultraviolet)

**Repository:** https://github.com/charmbracelet/ultraviolet
**Latest Version:** **No releases.** Pre-release / active development. No tagged version.
**License:** MIT
**Stars:** 350 | **Forks:** 36

### 2.1 Role in Nous

Ultraviolet is the low-level terminal primitives layer. It powers Bubbletea v2 and Lip Gloss v2 internally. Nous does NOT interact with Ultraviolet directly -- it is a transitive dependency consumed through Bubbletea and Lip Gloss. Understanding its architecture helps with debugging rendering issues and understanding performance characteristics.

### 2.2 Architecture

Ultraviolet is organized as layered primitives:

```
Terminal (lifecycle: raw mode, event loop, start/stop)
   |
   +-- TerminalScreen (screen state, alt screen, cursor, mouse, rendering)
         |
         +-- Screen interface (Bounds, CellAt, SetCell, WidthMethod)
               |
               +-- Buffer / Window (off-screen cell buffers)
               +-- screen package (drawing helpers: Context, Clear, Fill)
               +-- layout package (Cassowary constraint solver)
```

#### Layer Descriptions

| Layer | Type | Purpose |
|-------|------|---------|
| **Terminal** | `uv.DefaultTerminal()` | Manages raw mode, input event loop, start/stop lifecycle |
| **TerminalScreen** | `terminal.Screen()` | Screen state manager: rendering, alt screen, cursor, mouse modes, keyboard enhancements, bracketed paste, window title |
| **Screen** | Interface | Minimal: `Bounds()`, `CellAt()`, `SetCell()`, `WidthMethod()`. Implemented by TerminalScreen, Buffer, Window, ScreenBuffer |
| **Buffer** | Type | Flat grid of cells, off-screen. Implements Screen and Drawable |
| **Window** | Type | Parent/child relationships, shared-buffer views. Implements Screen and Drawable |
| **screen package** | Package | Drawing helpers: `Context` (styled text rendering), `Clear`, `Fill`, `Clone` |
| **layout package** | Package | Cassowary constraint-based layout solver: `Len`, `Min`, `Max`, `Percent`, `Ratio`, `Fill` |

### 2.3 Features

- **Cell-based diffing renderer** -- only redraws changed cells. Uses ECH/REP/ICH/DCH when available. Optimizes cursor movement. Minimal bandwidth, critical for SSH sessions.
- **Universal input** -- unified keyboard and mouse across platforms. Legacy encodings, Kitty keyboard protocol, SGR mouse, Windows Console input.
- **Inline and fullscreen** -- both alternate screen and inline modes. Inline preserves scrollback.
- **Cross-platform** -- Unix (termios + ANSI) and Windows (Console API).
- **Suspend/resume** -- `Stop()` and `Start()` for suspend/resume cycles, shelling out to editors.

### 2.4 Quick Start (Standalone)

```go
package main

import (
    "log"
    uv "github.com/charmbracelet/ultraviolet"
    "github.com/charmbracelet/ultraviolet/screen"
)

func main() {
    t := uv.DefaultTerminal()
    scr := t.Screen()
    scr.EnterAltScreen()

    if err := t.Start(); err != nil {
        log.Fatalf("failed to start: %v", err)
    }
    defer t.Stop()

    ctx := screen.NewContext(scr)
    for ev := range t.Events() {
        switch ev := ev.(type) {
        case uv.WindowSizeEvent:
            scr.Resize(ev.Width, ev.Height)
        case uv.KeyPressEvent:
            if ev.MatchString("q", "ctrl+c") {
                return
            }
        }
    }
}
```

### 2.5 Stability Assessment

- **No releases.** The repository has zero tagged releases as of June 2026.
- **"API may change"** -- explicitly stated in README.
- **Active development** -- 3 open issues, 13 open PRs.
- **Powers Bubbletea v2** -- consumed internally by the Charm ecosystem.
- **Pre-v0 / unstable.** The API is not semver-stable.

**Risk Assessment for Nous:**
- **Direct usage: HIGH risk.** Nous should NOT import Ultraviolet directly. No version pin is possible (no releases).
- **Transitive usage: LOW risk.** Bubbletea v2 and Lip Gloss v2 depend on it. Their pinned versions will pull a compatible Ultraviolet commit.
- **Recommendation:** Treat Ultraviolet as an internal implementation detail of Bubbletea/Lip Gloss. Do not write code against its API. If a need arises for low-level terminal control, encapsulate the dependency behind an interface.

---

## 3. Colorprofile (charmbracelet/colorprofile)

**Repository:** https://github.com/charmbracelet/colorprofile
**GoDoc:** https://pkg.go.dev/github.com/charmbracelet/colorprofile
**Latest Version:** `v0.3.2` (released 2025-08-13)
**License:** MIT
**Importers:** 21

### 3.1 Role in Nous

Colorprofile detects the terminal's color capabilities and downsamples colors automatically. It ensures Nous's rich styled output degrades gracefully across terminals -- from true-color terminals down to ASCII-only environments (e.g., piped output, CI logs).

### 3.2 Core Types

#### `Profile` type

```go
type Profile byte

const (
    NoTTY    Profile = iota  // No terminal
    Ascii                     // No color (monospace text only)
    ANSI                      // 16 colors (4-bit)
    ANSI256                   // 256 colors (8-bit)
    TrueColor                 // 16 million colors (24-bit)
)
```

#### Detection Functions

```go
// Detect from output writer + environment (recommended)
func Detect(output io.Writer, env []string) Profile

// Detect from environment variables only
func Env(env []string) Profile

// Detect from terminfo database
func Terminfo(term string) Profile

// Detect in tmux session
func Tmux(env []string) Profile
```

Detection respects: `NO_COLOR`, `CLICOLOR`, `CLICOLOR_FORCE`, `COLORTERM`, `TERM`.

Rules:
- `TERM=dumb` -> NoTTY (unless `CLICOLOR_FORCE=1`)
- `COLORTERM=truecolor` -> TrueColor upgrade
- `TERM=xterm-256color` -> ANSI256
- `NO_COLOR` takes precedence over `CLICOLOR`/`CLICOLOR_FORCE`
- `NO_COLOR` disables colors but NOT text decoration (bold, italic, etc.)

#### Color Conversion

```go
// Convert a color to the profile's supported range
converted := profile.Convert(color.RGBA{0x6b, 0x50, 0xff, 0xff})

// Manual conversion
ansi256Color := colorprofile.ANSI256.Convert(c)
ansiColor := colorprofile.ANSI.Convert(c)
```

#### `Writer` type

```go
// Automatically downsample ANSI output for the terminal
w := colorprofile.NewWriter(os.Stdout, os.Environ())
fmt.Fprintf(w, "\x1b[38;2;107;80;255mFancy!\x1b[m")

// Override profile
w.Profile = colorprofile.ANSI  // Force 4-bit
w.Profile = colorprofile.Ascii // Strip all color
w.Profile = colorprofile.NoTTY // Strip all ANSI
```

### 3.3 Usage Pattern for Nous

```go
import "github.com/charmbracelet/colorprofile"

// Detect once at startup
profile := colorprofile.Detect(os.Stderr, os.Environ())

// Use Writer for output that needs automatic downsampling
output := colorprofile.NewWriter(os.Stderr, os.Environ())
```

Lip Gloss and Bubbletea use Colorprofile internally. Nous typically does NOT need to use it directly unless writing raw ANSI output outside the Charm ecosystem.

### 3.4 Stability Assessment

- **Pre-v1 (v0.3.2)** -- not semver-stable, but API surface is tiny and mature.
- **21 importers** -- primarily used internally by Charm libraries.
- **Simple, focused** -- does one thing well. Low risk of breaking changes.
- **Risk:** Very low. Even if API changes, the package is small enough to vendor or replace.

---

## 4. Pinned Versions

| Package | Import Path | Version | Date |
|---------|------------|---------|------|
| Bubbletea | `charm.land/bubbletea/v2` | `v2.0.7` | 2026-06-01 |
| Ultraviolet | `github.com/charmbracelet/ultraviolet` | **No release** (commit hash) | N/A |
| Colorprofile | `github.com/charmbracelet/colorprofile` | `v0.3.2` | 2025-08-13 |

### go.mod directive

```
require (
    charm.land/bubbletea/v2 v2.0.7
    github.com/charmbracelet/colorprofile v0.3.2
    // Ultraviolet: transitive dependency via bubbletea/v2, do NOT pin directly
)
```

---

## 5. Integration Architecture

```
Nous CLI (Cobra + Fang)
       |
       +-- Command: "status"
       |     RunE -> launch Bubbletea program
       |
       v
  Bubbletea v2 (Elm Architecture)
       |
       +-- Model: Nous state engine
       +-- Update: handle user input, state transitions
       +-- View: render with Lip Gloss styles
       |
       +-- Internal: Ultraviolet (terminal primitives)
       +-- Internal: Colorprofile (color detection)
       |
       v
  Terminal output (cell-based diffing renderer)
```

### Key Integration Points

1. **Bubbletea is launched from Cobra `RunE`** -- not from Fang. Fang handles help/errors. Interactive screens are separate Bubbletea programs.
2. **One Bubbletea program per interactive command** -- each command with a TUI creates `tea.NewProgram(model).Run()`.
3. **Alt screen for full-screen UIs** -- set `view.AltScreen = true` in `View()` for dashboards, viewers. Inline mode for simple prompts.
4. **Mouse support** -- set `view.MouseMode = tea.MouseModeCellMotion` in `View()` for dashboards, `tea.MouseModeAllMotion` for hover interactions.
5. **Context cancellation** -- pass `tea.WithContext(ctx)` to tie Bubbletea lifetime to Cobra's context.
6. **Ultraviolet is implicit** -- no direct import needed.

---

## 6. Gaps & Considerations

1. **Bubbletea v2 import path** -- Uses `charm.land/bubbletea/v2` vanity domain. Verify go.sum resolves correctly at implementation time.
2. **Ultraviolet has no releases** -- cannot pin. Must rely on Bubbletea's go.sum to lock the transitive dependency. This is a risk if Ultraviolet introduces breaking changes between Bubbletea patch versions.
3. **View() returns tea.View** -- The v2 API uses `tea.View` (not `string`). Create with `tea.NewView("content")` or set content via `v.SetContent("content")`.
4. **v1 to v2 migration** -- if any third-party Bubbles components still target v1, they will need updates. Verify Bubbles compatibility in Phase 3.
5. **Performance** -- Ultraviolet's diffing renderer is optimized for SSH. This is critical for Nous's remote usage scenarios.
6. **Testing Bubbletea** -- use `tea.NewProgram(model, tea.WithInput(strings.NewReader(...)), tea.WithOutput(io.Discard))` for unit testing without a real terminal.

---

## References

- Bubbletea GitHub: https://github.com/charmbracelet/bubbletea
- Bubbletea GoDoc: https://pkg.go.dev/github.com/charmbracelet/bubbletea
- Ultraviolet GitHub: https://github.com/charmbracelet/ultraviolet
- Colorprofile GoDoc: https://pkg.go.dev/github.com/charmbracelet/colorprofile
- Elm Architecture: https://guide.elm-lang.org/architecture/
