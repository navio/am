package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SyncConfig holds remote sync settings for the am CLI.
type SyncConfig struct {
	Provider     string `json:"provider,omitempty"`
	GistID       string `json:"gist_id,omitempty"`
	GistFilename string `json:"gist_filename,omitempty"`
	LastSynced   string `json:"last_synced,omitempty"`
	LastHash     string `json:"last_hash,omitempty"`
}

// Config is the on-disk configuration file at ~/.config/am/config.json.
type Config struct {
	Sync *SyncConfig `json:"sync,omitempty"`
}

// Dir returns the configuration directory, honoring XDG_CONFIG_HOME.
func Dir() (string, error) {
	if v := os.Getenv("AM_CONFIG_DIR"); v != "" {
		return v, nil
	}
	if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
		return filepath.Join(v, "am"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "am"), nil
}

// Path returns the absolute path of the config file.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load reads the configuration file. A missing file returns an empty Config.
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := &Config{}
	if len(data) == 0 {
		return cfg, nil
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}

// Save writes the configuration file, creating the directory if necessary.
func Save(cfg *Config) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path, err := Path()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}
