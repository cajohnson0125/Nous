# CLI Foundation: Cobra & Fang

> Research document for M001 — Charmbracelet UI Stack Research
> Date: 2026-06-06
> Status: Phase 1

## Summary

This document covers the two packages forming Nous's CLI foundation: **Cobra** (spf13/cobra) for command structure and flag parsing, and **Fang** (charmbracelet/fang) for styled help output and batteries-included Cobra integration.

---

## 1. Cobra (spf13/cobra)

**Repository:** https://github.com/spf13/cobra
**GoDoc:** https://pkg.go.dev/github.com/spf13/cobra
**Latest Version:** `v1.10.2` (released 2025-12-04)
**License:** Apache 2.0
**Stars:** 44.1k | **Importers:** 260k+

### 1.1 Role in Nous

Cobra is the command framework. It provides the command tree, flag parsing, subcommand routing, help generation, shell completions, and man page generation. Every user-facing CLI interaction starts with a Cobra command.

### 1.2 Core Types

#### `Command` struct

The central type. Every CLI action is a `*cobra.Command`.

```go
type Command struct {
    Use   string   // Usage string: "command [flags] [args]"
    Short string   // Short description shown in help listing
    Long  string   // Long description shown in full help
    Example string // Example usage strings

    Aliases    []string
    SuggestFor []string
    GroupID    string // Group for subcommand grouping (v1.6.0+)

    Args PositionalArgs // Validator for positional arguments

    // Lifecycle hooks (executed in order):
    PersistentPreRun  func(cmd *Command, args []string)
    PersistentPreRunE func(cmd *Command, args []string) error
    PreRun            func(cmd *Command, args []string)
    PreRunE           func(cmd *Command, args []string) error
    Run               func(cmd *Command, args []string)
    RunE              func(cmd *Command, args []string) error
    PostRun           func(cmd *Command, args []string)
    PostRunE          func(cmd *Command, args []string) error
    PersistentPostRun func(cmd *Command, args []string)
    PersistentPostRunE func(cmd *Command, args []string) error

    Version string // Enables --version flag when non-empty

    SilenceErrors    bool // Suppress error printing
    SilenceUsage     bool // Suppress usage on error
    Hidden           bool // Hide from help
    TraverseChildren bool // Parse flags on parent commands

    ValidArgs         []Completion
    ValidArgsFunction CompletionFunc

    Annotations      map[string]string
    Deprecated       string
    CompletionOptions CompletionOptions
    FParseErrWhitelist FParseErrWhitelist
}
```

### 1.3 Key Methods on `Command`

| Method | Purpose |
|--------|---------|
| `Execute() error` | Run the command tree with `os.Args[1:]` |
| `ExecuteC() (*Command, error)` | Execute and return the resolved command |
| `ExecuteContext(ctx context.Context) error` | Execute with context |
| `AddCommand(cmds ...*Command)` | Add subcommands |
| `AddGroup(groups ...*Group)` | Group subcommands in help (v1.6.0+) |
| `Flags() *pflag.FlagSet` | Local flags |
| `PersistentFlags() *pflag.FlagSet` | Flags inherited by children |
| `SetContext(ctx context.Context)` | Set context |
| `Context() context.Context` | Get context |
| `SetOut(w io.Writer)` | Set stdout |
| `SetErr(w io.Writer)` | Set stderr |
| `SetHelpFunc(f)` | Custom help handler |
| `SetUsageFunc(f)` | Custom usage handler |
| `RegisterFlagCompletionFunc(name, f)` | Shell completion for a flag |
| `MarkFlagRequired(name)` | Require a flag |
| `MarkFlagsMutuallyExclusive(names)` | XOR flag group |
| `MarkFlagsRequiredTogether(names)` | AND flag group |
| `ValidateArgs(args)` | Validate positional args |
| `Find(args) (*Command, []string, error)` | Resolve command from args |

### 1.4 Argument Validators (PositionalArgs)

```go
NoArgs                  // Reject any args
ArbitraryArgs           // Allow any args
ExactArgs(n int)        // Require exactly n
MaximumNArgs(n int)     // At most n
MinimumNArgs(n int)     // At least n
RangeArgs(min, max int) // Between min and max
MatchAll(pargs ...)     // Combine validators
```

### 1.5 Flag Types (via pflag)

Cobra uses `pflag` (POSIX-compliant flag parsing). Flags are registered on a `*pflag.FlagSet`:

```go
flags := cmd.Flags()
flags.StringP("output", "o", "", "Output file")
flags.BoolP("verbose", "v", false, "Verbose output")
flags.IntP("port", "p", 8080, "Port number")
flags.StringSlice("tags", []string{}, "Tags")
flags.Duration("timeout", 30*time.Second, "Timeout")
```

### 1.6 Lifecycle Hook Execution Order

For a command with parent chain `root -> parent -> child`:

1. `root.PersistentPreRun`
2. `parent.PersistentPreRun`
3. `child.PersistentPreRun`
4. `child.PreRun`
5. `child.Run` (or `RunE`)
6. `child.PostRun`
7. `child.PersistentPostRun`
8. `parent.PersistentPostRun`
9. `root.PersistentPostRun`

When `EnableTraverseRunHooks = true`, all parent hooks execute. Otherwise, only the first found hook runs.

### 1.7 Shell Completions

Cobra generates completions for: bash, zsh, fish, powershell.

```go
// Register dynamic completion for a flag
cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
    return []cobra.Completion{"json", "yaml", "table"}, cobra.ShellCompDirectiveNoFileComp
})
```

`ShellCompDirective` values:
- `ShellCompDirectiveDefault` — default file completion
- `ShellCompDirectiveNoFileComp` — no file completion
- `ShellCompDirectiveNoSpace` — no space after completion
- `ShellCompDirectiveFilterFileExt` — filter by file extension
- `ShellCompDirectiveFilterDirs` — filter directories

### 1.8 Minimal Working Example

```go
package main

import (
    "fmt"
    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "nous",
        Short: "Nous — dynamic state engine",
        Long:  "Nous is a multi-layer, dynamic state engine for terminal workflows.",
        Args:  cobra.NoArgs,
    }

    applyCmd := &cobra.Command{
        Use:   "apply [path]",
        Short: "Apply a state configuration",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Printf("Applying: %s\n", args[0])
            return nil
        },
    }

    rootCmd.AddCommand(applyCmd)

    if err := rootCmd.Execute(); err != nil {
        // cobra already printed the error
    }
}
```

### 1.9 Stability Assessment

- **Mature and stable.** Used by Kubernetes, Hugo, GitHub CLI, and 260k+ packages.
- **Semver v1** — backward-compatible guarantees.
- **Active maintenance.** v1.10.2 released Dec 2025.
- **No breaking changes expected** in the v1.x line.

---

## 2. Fang (charmbracelet/fang)

**Repository:** https://github.com/charmbracelet/fang
**GoDoc:** https://pkg.go.dev/github.com/charmbracelet/fang
**Latest Version:** `v2.0.1` (released 2026-03-11)
**License:** MIT
**Stars:** 1.9k | **Importers:** 222

> **Note:** Fang v2.0.1 requires import path `charm.land/fang/v2`. The pkg.go.dev page for `fang` (no version suffix) shows v1.0.0. Use the v2 module path.

### 2.1 Role in Nous

Fang wraps Cobra to provide styled help pages, styled error output, automatic `--version`, man pages, and shell completions. It is the presentation layer on top of Cobra's command tree. In Nous, Fang ensures every CLI interaction looks polished with zero per-command effort.

### 2.2 Core API

#### `Execute` — the main entry point

```go
func Execute(ctx context.Context, root *cobra.Command, options ...Option) error
```

Replaces `cmd.Execute()`. Applies all Fang enhancements to the root command and executes it.

### 2.3 Options (Configuration)

| Option | Type | Purpose |
|--------|------|---------|
| `WithVersion(version string)` | Option | Set version string for `--version` |
| `WithCommit(commit string)` | Option | Set commit SHA (appended to version) |
| `WithColorSchemeFunc(cs ColorSchemeFunc)` | Option | Custom color scheme function (v0.2.0+) |
| `WithTheme(theme ColorScheme)` | Option | **Deprecated:** use WithColorSchemeFunc |
| `WithErrorHandler(handler ErrorHandler)` | Option | Custom error handler (v0.2.0+) |
| `WithNotifySignal(signals ...os.Signal)` | Option | Signals that interrupt execution (v0.3.0+) |
| `WithoutCompletions()` | Option | Disable `completion` subcommand |
| `WithoutManpage()` | Option | Disable `man` subcommand |
| `WithoutVersion()` | Option | Disable `--version` flag (v0.3.0+) |

### 2.4 ColorScheme

```go
type ColorScheme struct {
    Base           color.Color
    Title          color.Color
    Description    color.Color
    Codeblock      color.Color
    Program        color.Color
    DimmedArgument color.Color
    Comment        color.Color
    Flag           color.Color
    FlagDefault    color.Color
    Command        color.Color
    QuotedString   color.Color
    Argument       color.Color
    Help           color.Color
    Dash           color.Color
    ErrorHeader    [2]color.Color // 0=fg 1=bg
    ErrorDetails   color.Color
}
```

Built-in constructors:
- `DefaultColorScheme(c lipgloss.LightDarkFunc) ColorScheme` — adaptive light/dark (v0.2.0+)
- `AnsiColorScheme(c lipgloss.LightDarkFunc) ColorScheme` — ANSI color scheme (v0.2.0+)

### 2.5 Styles

```go
type Styles struct {
    Text            lipgloss.Style
    Title           lipgloss.Style
    Span            lipgloss.Style
    ErrorHeader     lipgloss.Style
    ErrorText       lipgloss.Style
    FlagDescription lipgloss.Style
    FlagDefault     lipgloss.Style
    Codeblock       Codeblock
    Program         Program
}

type Codeblock struct {
    Base    lipgloss.Style
    Program lipgloss.Style
    Text    lipgloss.Style
    Comment lipgloss.Style
}

type Program struct {
    Name           lipgloss.Style
    Command        lipgloss.Style
    Flag           lipgloss.Style
    Argument       lipgloss.Style
    DimmedArgument lipgloss.Style
    QuotedString   lipgloss.Style
}
```

### 2.6 ErrorHandler

```go
type ErrorHandler = func(w io.Writer, styles Styles, err error)
```

The default handler (`DefaultErrorHandler`) prints styled error output. Custom handlers receive the full `Styles` object for consistent theming.

### 2.7 What Fang Adds to Cobra

| Feature | How |
|---------|-----|
| Styled help pages | Replaces Cobra's default help template with Lip Gloss-styled output |
| Styled errors | Custom error handler with themed formatting |
| Automatic `--version` | From build info or explicit string |
| Man pages | Hidden `man` subcommand using mango (better than Cobra's roff) |
| Shell completions | `completion` subcommand auto-added |
| UX: Silent usage | Help is not shown after user errors (only on `--help`) |
| Light/dark theme | Adapts to terminal color scheme via `lipgloss.LightDarkFunc` |

### 2.8 Minimal Working Example

```go
package main

import (
    "context"
    "os"

    "charm.land/fang/v2"
    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "nous",
        Short: "Nous — dynamic state engine",
        Long:  "Nous is a multi-layer, dynamic state engine.",
    }

    applyCmd := &cobra.Command{
        Use:   "apply [path]",
        Short: "Apply a state configuration",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            return nil
        },
    }

    rootCmd.AddCommand(applyCmd)

    // Fang replaces rootCmd.Execute()
    if err := fang.Execute(context.Background(), rootCmd,
        fang.WithVersion("0.1.0"),
        fang.WithCommit("abc1234"),
    ); err != nil {
        os.Exit(1)
    }
}
```

### 2.9 Stability Assessment

- **Stable v2 release.** v2.0.1 released March 2026.
- **Experimental label** in README, but the API is clean and semver-stable at v2.
- **Active development** by Charm team (aymanbagabas, meowgorithm).
- **Import path changed** from v1 to v2: use `charm.land/fang/v2`.
- **Dependencies:** Fang pulls in Lip Gloss, Bubbletea, and other Charm libraries. This aligns with Nous's stack.
- **Risk:** Low. Even if Fang stops updating, it wraps Cobra in a thin layer that could be replaced with custom help templates.

---

## 3. Integration Pattern: Cobra + Fang in Nous

### 3.1 Architecture

```
User invokes CLI
       |
       v
   Fang.Execute(ctx, rootCmd, opts...)
       |
       |-- Injects styled help template
       |-- Injects styled error handler
       |-- Adds --version flag
       |-- Adds 'completion' subcommand
       |-- Adds hidden 'man' subcommand
       |
       v
   Cobra resolves command tree
       |
       |-- PersistentPreRun hooks
       |-- PreRun hooks
       |-- Run / RunE
       |-- PostRun hooks
       |
       v
   Fang renders styled output (help, errors)
```

### 3.2 Pattern for Nous Commands

```go
// cmd/root.go
package cmd

import (
    "context"
    "os"

    "charm.land/fang/v2"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "nous",
    Short: "Dynamic state engine",
    Long:  "Nous — a multi-layer, dynamic state engine for terminal workflows.",
}

func Execute() {
    if err := fang.Execute(context.Background(), rootCmd,
        fang.WithVersion("0.1.0"),
    ); err != nil {
        os.Exit(1)
    }
}

// cmd/apply.go
var applyCmd = &cobra.Command{
    Use:   "apply [path]",
    Short: "Apply a state configuration",
    Args:  cobra.ExactArgs(1),
    RunE:  runApply,
}

func init() {
    rootCmd.AddCommand(applyCmd)
    applyCmd.Flags().Bool("dry-run", false, "Preview changes without applying")
}
```

### 3.3 Key Integration Points

1. **Use `fang.Execute()` instead of `cobra.Execute()`** — single integration point.
2. **Do NOT call `cmd.SetHelpFunc()`** — Fang overrides it. If custom help is needed, set it before Fang's `Execute`.
3. **Version info** — pass via `WithVersion()` or let Fang detect from Go build info.
4. **Error handling** — Fang silences Cobra's default error printing and uses its own styled handler. Use `RunE` (not `Run`) for proper error propagation.
5. **Context** — Fang's `Execute` takes `context.Context`. Use it for cancellation and signal handling.

### 3.4 Command Groups (Cobra v1.6.0+)

Fang respects Cobra's `AddGroup()` for grouping subcommands in styled help output:

```go
rootCmd.AddGroup(
    &cobra.Group{ID: "core", Title: "Core Commands"},
    &cobra.Group{ID: "state", Title: "State Management"},
)
rootCmd.AddCommand(applyCmd)
applyCmd.GroupID = "core"
```

---

## 4. Pinned Versions

| Package | Import Path | Version | Date |
|---------|------------|---------|------|
| Cobra | `github.com/spf13/cobra` | `v1.10.2` | 2025-12-04 |
| Fang | `charm.land/fang/v2` | `v2.0.1` | 2026-03-11 |

### go.mod directive

```
require (
    github.com/spf13/cobra v1.10.2
    charm.land/fang/v2 v2.0.1
)
```

---

## 5. Gaps & Considerations

1. **Fang v2 module path** — Must use `charm.land/fang/v2` import path. The `pkg.go.dev/github.com/charmbracelet/fang` page shows v1.0.0; the v2 API uses the `charm.land` vanity domain.
2. **No Bubbletea integration** — Fang handles help/error rendering only. It does NOT launch a Bubbletea program. TUI interactions happen inside command `RunE` handlers.
3. **Custom error handling** — For domain-specific errors (validation, state conflicts), Nous should define a custom `ErrorHandler` that formats errors with Fang styles but adds domain context.
4. **Signal handling** — Fang's `WithNotifySignal()` handles OS signals. For Nous, pass `os.Interrupt, syscall.SIGTERM` for graceful shutdown of interactive TUI sessions.
5. **Testing** — Cobra commands are easily testable via `cmd.SetArgs()` and `cmd.Execute()`. When using Fang, test with `fang.Execute()` to match production behavior, or test commands directly with `cmd.Execute()` for unit tests.

---

## References

- Cobra GitHub: https://github.com/spf13/cobra
- Cobra GoDoc: https://pkg.go.dev/github.com/spf13/cobra
- Cobra.dev documentation: https://cobra.dev
- Fang GitHub: https://github.com/charmbracelet/fang
- Fang GoDoc (v1): https://pkg.go.dev/github.com/charmbracelet/fang
- Fang GoDoc (v2): https://pkg.go.dev/charm.land/fang/v2
- pflag: https://pkg.go.dev/github.com/spf13/pflag
