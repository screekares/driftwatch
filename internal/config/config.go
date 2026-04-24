package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the top-level driftwatch configuration.
type Config struct {
	Version  string            `json:"version"`
	Profiles map[string]Profile `json:"profiles"`
}

// Profile represents a named environment profile (e.g. staging, production).
type Profile struct {
	Name     string            `json:"name"`
	Provider string            `json:"provider"` // aws, gcp, azure
	Region   string            `json:"region"`
	Tags     map[string]string `json:"tags,omitempty"`
}

const defaultConfigFile = ".driftwatch.json"

// Load reads a Config from the given file path.
// If path is empty it falls back to .driftwatch.json in the working directory.
func Load(path string) (*Config, error) {
	if path == "" {
		path = defaultConfigFile
	}
	path = filepath.Clean(path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading %q: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parsing %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Version == "" {
		return fmt.Errorf("version field is required")
	}
	for name, p := range c.Profiles {
		if p.Provider == "" {
			return fmt.Errorf("profile %q: provider is required", name)
		}
	}
	return nil
}
