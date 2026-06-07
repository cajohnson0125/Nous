// Package internal contains the core application types for Nous.
package internal

import (
	"runtime/debug"
	"strings"
)

// App holds shared application state passed through commands.
type App struct {
	Name    string
	Version string
	Commit  string
}

// NewApp creates an App with version and commit info from build data.
// Falls back to "dev" if not set via ldflags.
func NewApp() *App {
	app := &App{
		Name:    "Nous",
		Version: "dev",
		Commit:  "unknown",
	}

	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range bi.Settings {
			switch setting.Key {
			case "vcs.revision":
				if len(setting.Value) >= 7 {
					app.Commit = setting.Value[:7]
				}
			case "vcs.modified":
				if setting.Value == "true" {
					app.Commit += "-dirty"
				}
			}
		}
		// Use the main module version only if it was explicitly set
		// (not a git pseudo-version like 0.0.0-20260607...).
		if bi.Main.Version != "" && !strings.HasPrefix(bi.Main.Version, "0.0.0-") {
			app.Version = strings.TrimPrefix(bi.Main.Version, "v")
		}
	}

	return app
}

// VersionString returns a display-formatted version string.
func (a *App) VersionString() string {
	if a.Version == "dev" && a.Commit == "unknown" {
		return "dev"
	}
	if a.Commit == "unknown" || a.Commit == "" {
		return a.Version
	}
	return a.Version + " (" + a.Commit + ")"
}
