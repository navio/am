package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ShellType represents the type of shell
type ShellType string

const (
	Bash    ShellType = "bash"
	Zsh     ShellType = "zsh"
	Fish    ShellType = "fish"
	Unknown ShellType = "unknown"
)

// Detector handles shell detection and dotfile path resolution
type Detector struct{}

// NewDetector creates a new Detector instance
func NewDetector() *Detector {
	return &Detector{}
}

// DetectShell determines the current shell from the SHELL environment variable
func (d *Detector) DetectShell() ShellType {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		return Unknown
	}

	// Extract the shell name from the path (e.g., /bin/zsh -> zsh)
	shellName := filepath.Base(shellPath)

	switch {
	case strings.Contains(shellName, "zsh"):
		return Zsh
	case strings.Contains(shellName, "bash"):
		return Bash
	case strings.Contains(shellName, "fish"):
		return Fish
	default:
		return Unknown
	}
}

// GetDotfilePath returns the path to the appropriate dotfile based on shell type
func (d *Detector) GetDotfilePath() (string, error) {
	shellType := d.DetectShell()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	var dotfilePath string
	switch shellType {
	case Zsh:
		dotfilePath = filepath.Join(homeDir, ".zshrc")
	case Bash:
		dotfilePath = filepath.Join(homeDir, ".bashrc")
	case Fish:
		dotfilePath = filepath.Join(homeDir, ".config", "fish", "config.fish")
	case Unknown:
		return "", fmt.Errorf("unable to detect shell type (SHELL=%s). Supported shells: bash, zsh, fish", os.Getenv("SHELL"))
	}

	return dotfilePath, nil
}

// GetShellType returns the current shell type
func (d *Detector) GetShellType() (ShellType, error) {
	shellType := d.DetectShell()
	if shellType == Unknown {
		return Unknown, fmt.Errorf("unable to detect shell type (SHELL=%s). Supported shells: bash, zsh, fish", os.Getenv("SHELL"))
	}
	return shellType, nil
}
