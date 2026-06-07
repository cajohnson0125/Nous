// Package config provides YAML configuration loading, saving, and path
// resolution for Nous. Config files are searched in priority order:
//
//  1. .nous.yaml  (project-local, highest priority)
//  2. nous.yaml   (project-local)
//  3. $XDG_CONFIG_HOME/nous/nous.yaml (user-level)
//
// If no config file is found, sensible defaults are returned.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

// Config represents the full Nous configuration.
// Fields must remain extensible — adding new fields must not break
// existing config files (yaml.v3 ignores unknown keys by default).
type Config struct {
	Theme  string `yaml:"theme"`
	Cursor string `yaml:"cursor"`
	Blink  bool   `yaml:"blink"`

	LLM LLMConfig `yaml:"llm"`
}

// LLMConfig holds provider-specific LLM settings.
// These are loaded but not used in M02 — wired up in M03.
type LLMConfig struct {
	Provider        string `yaml:"provider"`
	APIKey          string `yaml:"api_key"`
	BaseURL         string `yaml:"base_url"`
	Model           string `yaml:"model"`
	ReasoningEffort string `yaml:"reasoning_effort"`
	ThinkingBudget  int    `yaml:"thinking_budget"`
}

// projectLocalConfigFiles lists project-local config filenames
// in priority order (highest first).
var projectLocalConfigFiles = []string{
	".nous.yaml",
	"nous.yaml",
}

// xdgConfigPath is the XDG-relative config file path.
const xdgConfigPath = "nous/nous.yaml"

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Theme:  "dark",
		Cursor: "bar",
		Blink:  true,
		LLM: LLMConfig{
			Provider:        "zai",
			APIKey:          "$ZAI_API_KEY",
			BaseURL:         "https://api.z.ai/api/coding/paas/v4",
			Model:           "glm-4.7",
			ReasoningEffort: "none",
			ThinkingBudget:  0,
		},
	}
}

// Paths returns the full list of config file paths in search priority order
// (highest priority first).
func Paths() []string {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	result := make([]string, 0, len(projectLocalConfigFiles)+1)

	// Project-local paths (highest priority).
	for _, name := range projectLocalConfigFiles {
		result = append(result, filepath.Join(wd, name))
	}

	// User-level XDG path (lowest priority).
	if path := xdgConfigFilePath(); path != "" {
		result = append(result, path)
	}

	return result
}

// Load searches config paths in priority order and returns the first file
// found. Environment variable references in string fields (e.g. $ZAI_API_KEY)
// are expanded using os.Getenv. If no config file exists, defaults are
// returned without error.
func Load() (*Config, error) {
	// Try project-local files first, then XDG path.
	for _, name := range projectLocalConfigFiles {
		path := filepath.Join(mustGetwd(), name)
		if _, err := os.Stat(path); err == nil {
			return loadFromFile(path)
		}
	}

	// Try XDG config path.
	if xdgPath := xdgConfigFilePath(); xdgPath != "" {
		if _, err := os.Stat(xdgPath); err == nil {
			return loadFromFile(xdgPath)
		}
	}

	// No config file found — return defaults silently.
	return Default(), nil
}

// loadFromFile reads and parses a YAML config file, then expands
// environment variable references.
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	cfg := Default()
	// Unmarshal on top of defaults so any fields not in the file keep
	// their default values. yaml.v3 ignores unknown fields by default.
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	expandEnvVars(cfg)

	return cfg, nil
}

// Save writes the config as YAML to the user-level XDG config path.
// Parent directories are created automatically. Returns the path written.
func Save(cfg *Config) (string, error) {
	path := xdgConfigFilePath()
	if path == "" {
		return "", fmt.Errorf("cannot determine XDG config path")
	}

	// Marshal the config struct. Environment variable references (e.g. $VAR)
	// are kept as-is in the saved file — they are expanded only at load time.
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal config: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create config directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("write config %s: %w", path, err)
	}

	return path, nil
}

// expandEnvVars resolves environment variable references in string config
// fields. A reference starts with '$' followed by the variable name.
// If the environment variable is not set, the field keeps the raw reference.
func expandEnvVars(cfg *Config) {
	cfg.LLM.APIKey = expandEnvRef(cfg.LLM.APIKey)
}

// expandEnvRef expands a single $VAR_NAME reference.
// If the string does not start with '$', it is returned unchanged.
// If the environment variable is unset, the raw reference is returned.
func expandEnvRef(s string) string {
	if strings.HasPrefix(s, "$") {
		name := s[1:]
		if val, ok := os.LookupEnv(name); ok {
			return val
		}
	}
	return s
}

// xdgConfigFilePath returns the user-level XDG config path.
// It reads XDG_CONFIG_HOME from the environment at call time so that
// tests which override it via t.Setenv see the updated value.
func xdgConfigFilePath() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = xdg.ConfigHome
	}
	if configHome == "" {
		return ""
	}
	return filepath.Join(configHome, xdgConfigPath)
}

// UserConfigPath returns the path where Save writes the user-level config.
// This is the primary XDG config location, regardless of whether the file
// currently exists.
func UserConfigPath() (string, error) {
	path := xdgConfigFilePath()
	if path == "" {
		return "", fmt.Errorf("cannot determine XDG config path")
	}
	return path, nil
}

// UserConfigExists returns true if the user-level config file already exists.
func UserConfigExists() (bool, string, error) {
	path, err := UserConfigPath()
	if err != nil {
		return false, "", err
	}
	_, err = os.Stat(path)
	if err == nil {
		return true, path, nil
	}
	if os.IsNotExist(err) {
		return false, path, nil
	}
	return false, path, err
}

// mustGetwd returns the current working directory or "." on error.
func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}
