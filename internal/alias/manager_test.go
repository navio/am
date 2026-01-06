package alias

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/navio/am/internal/dotfile"
	"github.com/navio/am/internal/shell"
)

// setupTestEnv creates a temporary dotfile for testing
func setupTestEnv(t *testing.T, shellType shell.ShellType) (string, func()) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "am-test-*")
	if err != nil {
		t.Fatal(err)
	}

	// Create dotfile
	var dotfileName string
	if shellType == shell.Zsh {
		dotfileName = ".zshrc"
	} else {
		dotfileName = ".bashrc"
	}
	dotfilePath := filepath.Join(tmpDir, dotfileName)

	// Write initial content
	handler := dotfile.NewHandler()
	initialLines := []string{
		"# Test configuration",
		"alias ll='ls -la'",
		"alias gs='git status'",
		"export PATH=/usr/bin",
	}
	if err := handler.Write(dotfilePath, initialLines); err != nil {
		t.Fatal(err)
	}

	// Override home directory for tests
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	// Setup shell environment
	originalShell := os.Getenv("SHELL")
	if shellType == shell.Zsh {
		os.Setenv("SHELL", "/bin/zsh")
	} else {
		os.Setenv("SHELL", "/bin/bash")
	}

	// Cleanup function
	cleanup := func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("SHELL", originalShell)
		os.RemoveAll(tmpDir)
	}

	return dotfilePath, cleanup
}

func TestLoadAliases(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()
	aliasFile, err := manager.LoadAliases()
	if err != nil {
		t.Fatalf("LoadAliases() error = %v", err)
	}

	if len(aliasFile.Aliases) != 2 {
		t.Errorf("LoadAliases() found %d aliases, want 2", len(aliasFile.Aliases))
	}

	// Check first alias
	if aliasFile.Aliases[0].Name != "ll" {
		t.Errorf("First alias name = %v, want 'll'", aliasFile.Aliases[0].Name)
	}
	if aliasFile.Aliases[0].Command != "ls -la" {
		t.Errorf("First alias command = %v, want 'ls -la'", aliasFile.Aliases[0].Command)
	}
}

func TestAdd(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	// Add a new alias
	err := manager.Add("serve", "python -m http.server 8000")
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	// Verify it was added
	aliases, err := manager.List("")
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, a := range aliases {
		if a.Name == "serve" && a.Command == "python -m http.server 8000" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Add() alias not found in list")
	}
}

func TestAddDuplicate(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	// Try to add a duplicate alias
	err := manager.Add("ll", "ls -lh")
	if err == nil {
		t.Error("Add() should error on duplicate alias")
	}
}

func TestAddInvalidName(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	tests := []struct {
		name      string
		aliasName string
		command   string
	}{
		{
			name:      "empty name",
			aliasName: "",
			command:   "ls -la",
		},
		{
			name:      "name with spaces",
			aliasName: "my alias",
			command:   "ls -la",
		},
		{
			name:      "name with special chars",
			aliasName: "my@alias",
			command:   "ls -la",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.Add(tt.aliasName, tt.command)
			if err == nil {
				t.Error("Add() should error on invalid alias name")
			}
		})
	}
}

func TestAddEmptyCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	err := manager.Add("test", "")
	if err == nil {
		t.Error("Add() should error on empty command")
	}
}

func TestList(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	aliases, err := manager.List("")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(aliases) != 2 {
		t.Errorf("List() returned %d aliases, want 2", len(aliases))
	}
}

func TestListWithSearch(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	// Search for "gs"
	aliases, err := manager.List("gs")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(aliases) != 1 {
		t.Errorf("List('gs') returned %d aliases, want 1", len(aliases))
	}

	if len(aliases) > 0 && aliases[0].Name != "gs" {
		t.Errorf("List('gs') returned alias %v, want 'gs'", aliases[0].Name)
	}

	// Search for non-existent
	aliases, err = manager.List("nonexistent")
	if err != nil {
		t.Fatal(err)
	}

	if len(aliases) != 0 {
		t.Errorf("List('nonexistent') returned %d aliases, want 0", len(aliases))
	}
}

func TestDelete(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	// Delete an existing alias
	err := manager.Delete("ll")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify it was deleted
	aliases, err := manager.List("")
	if err != nil {
		t.Fatal(err)
	}

	for _, a := range aliases {
		if a.Name == "ll" {
			t.Error("Delete() alias still exists after deletion")
		}
	}

	// Should have 1 alias remaining
	if len(aliases) != 1 {
		t.Errorf("After Delete(), got %d aliases, want 1", len(aliases))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	err := manager.Delete("nonexistent")
	if err == nil {
		t.Error("Delete() should error on non-existent alias")
	}
}

func TestUpdate(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	// Update an existing alias
	err := manager.Update("ll", "ls -lah")
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Verify it was updated
	alias, err := manager.Get("ll")
	if err != nil {
		t.Fatal(err)
	}

	if alias.Command != "ls -lah" {
		t.Errorf("Update() command = %v, want 'ls -lah'", alias.Command)
	}
}

func TestUpdateNonExistent(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	err := manager.Update("nonexistent", "ls -la")
	if err == nil {
		t.Error("Update() should error on non-existent alias")
	}
}

func TestUpdateEmptyCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	err := manager.Update("ll", "")
	if err == nil {
		t.Error("Update() should error on empty command")
	}
}

func TestGet(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	alias, err := manager.Get("ll")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if alias.Name != "ll" {
		t.Errorf("Get() name = %v, want 'll'", alias.Name)
	}

	if alias.Command != "ls -la" {
		t.Errorf("Get() command = %v, want 'ls -la'", alias.Command)
	}
}

func TestGetNonExistent(t *testing.T) {
	_, cleanup := setupTestEnv(t, shell.Zsh)
	defer cleanup()

	manager := NewManager()

	_, err := manager.Get("nonexistent")
	if err == nil {
		t.Error("Get() should error on non-existent alias")
	}
}

func TestValidateAliasName(t *testing.T) {
	manager := NewManager()

	validNames := []string{"ll", "gs", "my-alias", "my_alias", "alias123"}
	for _, name := range validNames {
		if err := manager.validateAliasName(name); err != nil {
			t.Errorf("validateAliasName(%v) should be valid, got error: %v", name, err)
		}
	}

	invalidNames := []string{"", "my alias", "my@alias", "my.alias", "my$alias"}
	for _, name := range invalidNames {
		if err := manager.validateAliasName(name); err == nil {
			t.Errorf("validateAliasName(%v) should be invalid", name)
		}
	}
}
