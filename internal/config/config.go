/*
Copyright © 2025 Lance Security <support@lancesecurity.org>
*/
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	EngagementRoot string `json:"engagement_root"`
	TemplateDir    string `json:"template_dir"`
}

// Default returns the configuration values
// and paths to use if no config file is found
func Default() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	cfgDir := filepath.Join(home, ".config", "praetor")

	return Config{
		EngagementRoot: filepath.Join(home, "engagements"),
		TemplateDir:    filepath.Join(cfgDir, "template"),
	}, nil
}

// Load reads the configuration from the config file, or returns defaults
// if the file does not exist
func Load() (Config, error) {
	cfg, err := Default()
	if err != nil {
		return Config{}, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	path := filepath.Join(home, ".config", "praetor", "config.json")

	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.EngagementRoot == "" {
		def, err := Default()
		if err != nil {
			return Config{}, err
		}
		cfg.EngagementRoot = def.EngagementRoot
	}
	if cfg.TemplateDir == "" {
		def, err := Default()
		if err != nil {
			return Config{}, err
		}
		cfg.TemplateDir = def.TemplateDir
	}

	return cfg, nil
}
