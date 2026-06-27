package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the tool configuration.
type Config struct {
	DefaultTemplate string            `yaml:"defaultTemplate"`
	DefaultBranch   string            `yaml:"defaultBranch"`
	Templates       map[string]string `yaml:"templates"`
}

// DefaultConfig returns a configuration with sensible default settings.
func DefaultConfig() *Config {
	return &Config{
		DefaultTemplate: "starter",
		DefaultBranch:   "main",
		Templates: map[string]string{
			"starter": "templates",
		},
	}
}

// GetConfigPath returns the path to the configuration file (~/.gitnew/config.yaml).
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gitnew", "config.yaml"), nil
}

// Load loads the configuration from the default path.
// If the file does not exist, it returns the default configuration.
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	return LoadFromFile(path)
}

// LoadFromFile loads the configuration from a specific file path.
func LoadFromFile(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	if cfg.DefaultBranch == "" {
		cfg.DefaultBranch = "main"
	}

	if cfg.Templates == nil {
		cfg.Templates = make(map[string]string)
	}

	if _, ok := cfg.Templates["starter"]; !ok {
		cfg.Templates["starter"] = "templates"
	}

	return cfg, nil
}
