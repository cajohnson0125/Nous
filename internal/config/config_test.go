package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cajohnson0125/Nous/internal/config"
)

// TestLoad_Defaults verifies that Load returns defaults when no config file
// exists anywhere (no error, no crash).
func TestLoad_Defaults(t *testing.T) {
	// Change to a temporary directory to ensure no project-local config exists.
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// Also ensure XDG config path doesn't find anything by setting
	// XDG_CONFIG_HOME to a temp dir.
	tmpXDG := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpXDG)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	def := config.Default()
	if cfg.Theme != def.Theme {
		t.Errorf("Theme = %q, want %q", cfg.Theme, def.Theme)
	}
	if cfg.Cursor != def.Cursor {
		t.Errorf("Cursor = %q, want %q", cfg.Cursor, def.Cursor)
	}
	if cfg.Blink != def.Blink {
		t.Errorf("Blink = %v, want %v", cfg.Blink, def.Blink)
	}
	if cfg.LLM.Provider != def.LLM.Provider {
		t.Errorf("LLM.Provider = %q, want %q", cfg.LLM.Provider, def.LLM.Provider)
	}
	if cfg.LLM.APIKey != def.LLM.APIKey {
		t.Errorf("LLM.APIKey = %q, want %q", cfg.LLM.APIKey, def.LLM.APIKey)
	}
	if cfg.LLM.BaseURL != def.LLM.BaseURL {
		t.Errorf("LLM.BaseURL = %q, want %q", cfg.LLM.BaseURL, def.LLM.BaseURL)
	}
	if cfg.LLM.Model != def.LLM.Model {
		t.Errorf("LLM.Model = %q, want %q", cfg.LLM.Model, def.LLM.Model)
	}
	if cfg.LLM.ThinkingBudget != def.LLM.ThinkingBudget {
		t.Errorf("LLM.ThinkingBudget = %d, want %d", cfg.LLM.ThinkingBudget, def.LLM.ThinkingBudget)
	}
}

// TestLoad_FromFile verifies that Load reads a YAML config file and
// overrides defaults with file values.
func TestLoad_FromFile(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// Also ensure XDG config path doesn't find anything.
	tmpXDG := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpXDG)

	// Write a project-local config file.
	yamlContent := `theme: light
cursor: block
blink: false
llm:
  provider: openai
  api_key: $MY_KEY
  base_url: https://api.openai.com/v1
  model: gpt-4
  reasoning_effort: high
  thinking_budget: 100
`
	configPath := filepath.Join(tmpDir, ".nous.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Theme != "light" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "light")
	}
	if cfg.Cursor != "block" {
		t.Errorf("Cursor = %q, want %q", cfg.Cursor, "block")
	}
	if cfg.Blink != false {
		t.Errorf("Blink = %v, want %v", cfg.Blink, false)
	}
	if cfg.LLM.Provider != "openai" {
		t.Errorf("LLM.Provider = %q, want %q", cfg.LLM.Provider, "openai")
	}
	if cfg.LLM.Model != "gpt-4" {
		t.Errorf("LLM.Model = %q, want %q", cfg.LLM.Model, "gpt-4")
	}
	if cfg.LLM.ThinkingBudget != 100 {
		t.Errorf("LLM.ThinkingBudget = %d, want %d", cfg.LLM.ThinkingBudget, 100)
	}
}

// TestSave verifies that Save writes a YAML config file to the XDG path
// and returns the written path.
func TestSave(t *testing.T) {
	tmpXDG := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpXDG)

	cfg := config.Default()
	cfg.Theme = "monokai"

	path, err := config.Save(cfg)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	if path == "" {
		t.Fatal("Save() returned empty path")
	}

	// Verify the file was created.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read saved file: %v", err)
	}

	if !strings.Contains(string(data), "theme: monokai") {
		t.Errorf("saved YAML does not contain 'theme: monokai', got:\n%s", string(data))
	}

	// Verify the path is under XDG_CONFIG_HOME.
	if !strings.HasPrefix(path, tmpXDG) {
		t.Errorf("saved path %q not under XDG_CONFIG_HOME %q", path, tmpXDG)
	}
}

// TestPaths_Priority verifies that Paths returns config paths in the
// correct priority order: .nous.yaml, nous.yaml, XDG path.
func TestPaths_Priority(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	tmpXDG := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpXDG)

	paths := config.Paths()

	if len(paths) < 3 {
		t.Fatalf("Paths() returned %d entries, want at least 3", len(paths))
	}

	// First should be .nous.yaml in cwd.
	want0 := filepath.Join(tmpDir, ".nous.yaml")
	if paths[0] != want0 {
		t.Errorf("Paths()[0] = %q, want %q", paths[0], want0)
	}

	// Second should be nous.yaml in cwd.
	want1 := filepath.Join(tmpDir, "nous.yaml")
	if paths[1] != want1 {
		t.Errorf("Paths()[1] = %q, want %q", paths[1], want1)
	}

	// Third should be the XDG user-level path.
	want2 := filepath.Join(tmpXDG, "nous", "nous.yaml")
	if paths[2] != want2 {
		t.Errorf("Paths()[2] = %q, want %q", paths[2], want2)
	}
}

// TestEnvVarExpansion verifies that $VAR references in the api_key field
// are resolved via os.Getenv at load time.
func TestEnvVarExpansion(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// Ensure XDG doesn't find anything.
	tmpXDG := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpXDG)

	// Set the env var that the config references.
	t.Setenv("TEST_NOUS_KEY", "resolved-secret-value")

	yamlContent := `llm:
  api_key: $TEST_NOUS_KEY
`
	configPath := filepath.Join(tmpDir, ".nous.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.LLM.APIKey != "resolved-secret-value" {
		t.Errorf("LLM.APIKey = %q, want %q", cfg.LLM.APIKey, "resolved-secret-value")
	}
}

// TestEnvVarExpansion_Unset verifies that an unset env var keeps the
// raw reference.
func TestEnvVarExpansion_Unset(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	tmpXDG := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpXDG)

	// Ensure the variable is not set.
	os.Unsetenv("NONEXISTENT_NOUS_VAR_12345")

	yamlContent := `llm:
  api_key: $NONEXISTENT_NOUS_VAR_12345
`
	configPath := filepath.Join(tmpDir, ".nous.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.LLM.APIKey != "$NONEXISTENT_NOUS_VAR_12345" {
		t.Errorf("LLM.APIKey = %q, want raw reference %q", cfg.LLM.APIKey, "$NONEXISTENT_NOUS_VAR_12345")
	}
}
