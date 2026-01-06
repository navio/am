package shell

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectShell(t *testing.T) {
	tests := []struct {
		name     string
		shellEnv string
		want     ShellType
	}{
		{
			name:     "zsh shell",
			shellEnv: "/bin/zsh",
			want:     Zsh,
		},
		{
			name:     "bash shell",
			shellEnv: "/bin/bash",
			want:     Bash,
		},
		{
			name:     "zsh with custom path",
			shellEnv: "/usr/local/bin/zsh",
			want:     Zsh,
		},
		{
			name:     "bash with custom path",
			shellEnv: "/usr/local/bin/bash",
			want:     Bash,
		},
		{
			name:     "unknown shell",
			shellEnv: "/bin/fish",
			want:     Unknown,
		},
		{
			name:     "empty shell",
			shellEnv: "",
			want:     Unknown,
		},
	}

	detector := NewDetector()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original SHELL
			originalShell := os.Getenv("SHELL")
			defer os.Setenv("SHELL", originalShell)

			// Set test SHELL
			os.Setenv("SHELL", tt.shellEnv)

			got := detector.DetectShell()
			if got != tt.want {
				t.Errorf("DetectShell() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDotfilePath(t *testing.T) {
	tests := []struct {
		name        string
		shellEnv    string
		wantName    string
		wantError   bool
	}{
		{
			name:      "zsh dotfile",
			shellEnv:  "/bin/zsh",
			wantName:  ".zshrc",
			wantError: false,
		},
		{
			name:      "bash dotfile",
			shellEnv:  "/bin/bash",
			wantName:  ".bashrc",
			wantError: false,
		},
		{
			name:      "unknown shell error",
			shellEnv:  "/bin/fish",
			wantError: true,
		},
	}

	detector := NewDetector()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original SHELL
			originalShell := os.Getenv("SHELL")
			defer os.Setenv("SHELL", originalShell)

			// Set test SHELL
			os.Setenv("SHELL", tt.shellEnv)

			got, err := detector.GetDotfilePath()
			if tt.wantError {
				if err == nil {
					t.Errorf("GetDotfilePath() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GetDotfilePath() unexpected error: %v", err)
				return
			}

			if filepath.Base(got) != tt.wantName {
				t.Errorf("GetDotfilePath() = %v, want filename %v", got, tt.wantName)
			}

			// Verify it's in the home directory
			homeDir, _ := os.UserHomeDir()
			expectedPath := filepath.Join(homeDir, tt.wantName)
			if got != expectedPath {
				t.Errorf("GetDotfilePath() = %v, want %v", got, expectedPath)
			}
		})
	}
}

func TestGetShellType(t *testing.T) {
	tests := []struct {
		name      string
		shellEnv  string
		want      ShellType
		wantError bool
	}{
		{
			name:      "valid zsh",
			shellEnv:  "/bin/zsh",
			want:      Zsh,
			wantError: false,
		},
		{
			name:      "valid bash",
			shellEnv:  "/bin/bash",
			want:      Bash,
			wantError: false,
		},
		{
			name:      "unknown shell",
			shellEnv:  "/bin/fish",
			wantError: true,
		},
	}

	detector := NewDetector()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original SHELL
			originalShell := os.Getenv("SHELL")
			defer os.Setenv("SHELL", originalShell)

			// Set test SHELL
			os.Setenv("SHELL", tt.shellEnv)

			got, err := detector.GetShellType()
			if tt.wantError {
				if err == nil {
					t.Errorf("GetShellType() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GetShellType() unexpected error: %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("GetShellType() = %v, want %v", got, tt.want)
			}
		})
	}
}
