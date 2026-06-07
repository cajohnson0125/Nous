# UI Stack Synthesis: Integration Mapping & v0.1 Package Set

> Research document for M001 -- Charmbracelet UI Stack Research
> Date: 2026-06-06
> Status: Phase 5

## Summary

This document consolidates all package research into a single integration map for Nous's terminal UI stack. It defines the minimal v0.1 package set, documents integration patterns, finalizes the Ultraviolet risk assessment, provides a data flow diagram, and identifies gaps for future work.

---

## 1. Package-to-Surface Mapping

Each package maps to a specific surface area in Nous:

| Package | Version | Import Path | Surface Area |
|---------|---------|------------|--------------|
| Cobra | v1.10.2 | `github.com/spf13/cobra` | Command tree, flag parsing, subcommand routing, help generation |
| Fang | v2.0.1 | `charm.land/fang/v2` | Styled help pages, styled errors, auto --version, manpages, completions |
| Bubbletea | v2.0.7 | `charm.land/bubbletea/v2` | Interactive TUI programs (Elm Architecture: Model-Update-View) |
| Lip Gloss | v2.0.3 | `charm.land/lipgloss/v2` | Declarative styling, layout, borders, colors, tables, lists, trees |
| Bubbles | v2.1.0 | `charm.land/bubbles/v2` | Pre-built TUI components: table, viewport, list, spinner, input, progress |
| Huh | v2.0.3 | `charm.land/huh/v2` | Interactive forms, prompts, approval gates, configuration wizards |
| Glamour | v2.0.0 | `charm.land/glamour/v2` | Markdown rendering with terminal styles, rich output display |
| Colorprofile | v0.3.2 | `github.com/charmbracelet/colorprofile` | Color detection and downsampling (transitive via Lip Gloss) |
| Ultraviolet | (no release) | `github.com/charmbracelet/ultraviolet` | Terminal primitives (transitive via Bubbletea, do NOT import directly) |

---

## 2. Minimal v0.1 Package Set

```
require (
    // CLI Foundation
    github.com/spf13/cobra v1.10.2
    charm.land/fang/v2 v2.0.1

    // TUI Framework
    charm.land/bubbletea/v2 v2.0.7
    github.com/charmbracelet/colorprofile v0.3.2
    // Ultraviolet: transitive via bubbletea/v2, do NOT pin directly

    // Styling & Layout
    charm.land/lipgloss/v2 v2.0.3
    charm.land/bubbles/v2 v2.1.0

    // User Interaction
    charm.land/huh/v2 v2.0.3
    charm.land/glamour/v2 v2.0.0
)
```

**Total direct dependencies: 8 packages.**
**Transitive dependency (do not import): Ultraviolet.**

---

## 3. Integration Patterns

### 3.1 Layer Architecture

```
Layer 1: CLI Entry Point
  Cobra (command routing) + Fang (styled help/errors)
       |
       v
Layer 2: Interactive TUI
  Bubbletea (Elm Architecture programs)
       |
       +-- Styling: Lip Gloss (all visual output)
       +-- Components: Bubbles (tables, lists, viewports, spinners)
       +-- Forms: Huh (prompts, confirmations, selects)
       +-- Rich Content: Glamour (markdown rendering)
       |
       v
Layer 3: Terminal Primitives (internal)
  Ultraviolet (diffing renderer, input handling)
  Colorprofile (color detection)
```

### 3.2 Pattern A: Simple CLI Command (non-interactive)

```
Cobra RunE -> print output with Lip Gloss styles
```

No Bubbletea needed. Just Cobra + Fang + Lip Gloss.

### 3.3 Pattern B: Interactive TUI Screen

```
Cobra RunE -> tea.NewProgram(model).Run()
  Model.View() uses Lip Gloss styles + Bubbles components
```

### 3.4 Pattern C: Form / Approval Gate

```
Cobra RunE -> huh.NewForm(groups...).Run()
  OR: Embed Huh fields in Bubbletea model
```

### 3.5 Pattern D: Rich Document Viewer

```
Cobra RunE -> tea.NewProgram(viewerModel).Run()
  viewerModel.View() uses Bubbles viewport + Glamour rendering
```

### 3.6 Pattern E: Dashboard

```
Cobra RunE -> tea.NewProgram(dashModel).Run()
  dashModel.View() uses Lip Gloss layout + Bubbles table + spinner
  Layout: JoinHorizontal/JoinVertical for panes
```

---

## 4. Data Flow Diagram

```
                        User invokes CLI
                              |
                              v
                    fang.Execute(ctx, rootCmd)
                         |
                         |-- Fang injects styled help/error
                         |-- Fang adds --version, completion, man
                         |
                         v
                    Cobra resolves command
                         |
              +----------+-----------+
              |          |           |
              v          v           v
         Pattern A   Pattern B   Pattern C
         (Simple)    (TUI)      (Form)
              |          |           |
              v          v           v
         Lip Gloss   Bubbletea    Huh Form
         styled      Program      .Run()
         output           |
                     +----+----+
                     |    |    |
                     v    v    v
                  Lip   Bub  Gla-
                  Gloss bles mour
                  (style)(comp)(md)
                     |    |    |
                     +----+----+
                          |
                          v
                   View() returns tea.View
                          |
                          v
                   Ultraviolet renderer
                   (diffing, cell-based)
                          |
                          v
                   Terminal output
```

---

## 5. Ultraviolet Risk Assessment (Finalized)

| Factor | Assessment |
|--------|-----------|
| **Stability** | No releases, pre-v0. API may change |
| **Exposure** | Nous does NOT import Ultraviolet directly |
| **Dependency chain** | Bubbletea v2 -> Ultraviolet (commit hash in go.sum) |
| **Risk of breakage** | LOW -- Bubbletea pins a specific commit |
| **Mitigation** | If Ultraviolet breaks, pin Bubbletea to a known-good version |
| **Recommendation** | Do NOT import directly. Treat as internal Bubbletea detail. |

**Final verdict: ACCEPTABLE RISK.** Ultraviolet's instability is contained by Bubbletea's semver guarantees. If a Bubbletea patch release pulls a broken Ultraviolet commit, pin the Bubbletea version.

---

## 6. Gaps Identified

### 6.1 Import Path Ambiguities

Several packages use different import paths between v1 and v2:

| Package | v1 Import | v2 Import | Notes |
|---------|-----------|-----------|-------|
| Bubbletea | `github.com/charmbracelet/bubbletea` | `charm.land/bubbletea/v2` | VANITY DOMAIN |
| Lip Gloss | `github.com/charmbracelet/lipgloss` | `charm.land/lipgloss/v2` | VANITY DOMAIN |
| Bubbles | `github.com/charmbracelet/bubbles` | `charm.land/bubbles/v2` | VANITY DOMAIN |
| Huh | `github.com/charmbracelet/huh` | `charm.land/huh/v2` | VANITY DOMAIN |
| Fang | `github.com/charmbracelet/fang` | `charm.land/fang/v2` | VANITY DOMAIN |
| Glamour | `github.com/charmbracelet/glamour` | `charm.land/glamour/v2` | VANITY DOMAIN |

**Action required:** At implementation time, run `go get` for each package and verify the `charm.land` vanity domains resolve correctly in go.sum.

### 6.2 Missing Features

1. **No horizontal scrolling in Bubbles table** -- For wide datasets, consider viewport wrapping or custom implementation.
2. **No built-in chart/graph rendering** -- Charm ecosystem has no chart package. For status dashboards with graphs, use Unicode block characters or ASCII art.
3. **No file watching** -- For `nous watch` functionality, use `fsnotify` separately from the Charm stack.
4. **No SSH server integration** -- For remote Nous sessions, use `github.com/charmbracelet/ssh` (Wish) separately.
5. **No persistent config** -- Fang and Cobra handle CLI config but not persistent state. Use Viper or a custom solution.

### 6.3 Version Compatibility Matrix

All packages in the v0.1 set use Charm v2 ecosystem APIs:

| Package Pair | Compatible | Notes |
|-------------|-----------|-------|
| Bubbletea v2 + Bubbles v2 | YES | Same v2 ecosystem |
| Bubbletea v2 + Lip Gloss v2 | YES | Bubbletea uses Lip Gloss internally |
| Bubbletea v2 + Huh v2 | YES | Huh is a Bubbletea program |
| Lip Gloss v2 + Glamour v2 | YES | Glamour uses Lip Gloss v2 |
| Fang v2 + Cobra v1 | YES | Fang wraps Cobra v1.x |
| Fang v2 + Lip Gloss v2 | YES | Fang uses Lip Gloss for styling |

**No compatibility conflicts detected.**

---

## 7. Recommendations for v0.1

1. **Start with Cobra + Fang** -- Wire up the command tree and get styled help working. This is the foundation.
2. **Add Lip Gloss** -- Define the Nous color palette and base styles as a shared package.
3. **Build one Bubbletea screen** -- Pick the simplest interactive screen (e.g., `nous status --watch`) and implement it with Lip Gloss + Bubbles.
4. **Add Huh for approvals** -- Implement the approval gate flow with Huh Confirm.
5. **Add Glamour for output** -- Render Markdown reports in viewports.
6. **Test color degradation** -- Verify all output degrades gracefully with `NO_COLOR=1` and `TERM=dumb`.

---

## References

- Phase 1: `docs/project/research/cli-foundation.md`
- Phase 2: `docs/project/research/tui-framework.md`
- Phase 3: `docs/project/research/styling-layout.md`
- Phase 4: `docs/project/research/user-interaction.md`
- Charm ecosystem: https://charm.sh
- Cobra.dev: https://cobra.dev
