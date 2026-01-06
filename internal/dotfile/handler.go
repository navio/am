package dotfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/navio/alias-manager/internal/shell"
)

// Handler manages reading, writing, and backing up dotfiles
type Handler struct{}

// NewHandler creates a new Handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// Read reads the entire dotfile and returns its lines
func (h *Handler) Read(path string) ([]string, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist, return empty lines (will be created on first write)
		return []string{}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return lines, nil
}

// Write writes lines to the dotfile, creating a backup first if the file exists
func (h *Handler) Write(path string, lines []string) error {
	// Create backup if file exists
	if _, err := os.Stat(path); err == nil {
		if err := h.Backup(path); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Write to a temporary file first
	tmpPath := path + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	writer := bufio.NewWriter(file)
	for i, line := range lines {
		if _, err := writer.WriteString(line); err != nil {
			file.Close()
			os.Remove(tmpPath)
			return fmt.Errorf("failed to write line %d: %w", i, err)
		}
		if _, err := writer.WriteString("\n"); err != nil {
			file.Close()
			os.Remove(tmpPath)
			return fmt.Errorf("failed to write newline at line %d: %w", i, err)
		}
	}

	if err := writer.Flush(); err != nil {
		file.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	if err := file.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Rename temporary file to actual file (atomic operation on Unix)
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// Backup creates a backup of the dotfile with a .bak extension
func (h *Handler) Backup(path string) error {
	backupPath := path + ".bak"

	src, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// Create creates a new dotfile with appropriate header comments
func (h *Handler) Create(path string, shellType shell.ShellType) error {
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return nil // File already exists, nothing to do
	}

	var header string
	switch shellType {
	case shell.Zsh:
		header = "# Zsh configuration file\n# Managed by alias-manager\n\n"
	case shell.Bash:
		header = "# Bash configuration file\n# Managed by alias-manager\n\n"
	default:
		header = "# Shell configuration file\n# Managed by alias-manager\n\n"
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	return nil
}

// EnsureExists ensures the dotfile exists, creating it if necessary
func (h *Handler) EnsureExists(path string, shellType shell.ShellType) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return h.Create(path, shellType)
	}
	return nil
}

// ReadOrCreate reads the dotfile or creates it if it doesn't exist
func (h *Handler) ReadOrCreate(path string, shellType shell.ShellType) ([]string, error) {
	if err := h.EnsureExists(path, shellType); err != nil {
		return nil, err
	}

	lines, err := h.Read(path)
	if err != nil {
		return nil, err
	}

	// If the file was just created, it will have the header
	// If it's an existing file that was empty, we might want to add a header
	if len(lines) == 0 {
		var header []string
		switch shellType {
		case shell.Zsh:
			header = []string{"# Zsh configuration file", "# Managed by alias-manager", ""}
		case shell.Bash:
			header = []string{"# Bash configuration file", "# Managed by alias-manager", ""}
		default:
			header = []string{"# Shell configuration file", "# Managed by alias-manager", ""}
		}
		return header, nil
	}

	// Ensure file ends with a newline for proper formatting
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) != "" {
		// Add a blank line at the end if the last line is not empty
		lines = append(lines, "")
	}

	return lines, nil
}
