package dotfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/navio/am/internal/shell"
)

func TestReadNonExistentFile(t *testing.T) {
	handler := NewHandler()

	// Try to read a file that doesn't exist
	lines, err := handler.Read("/tmp/nonexistent-file-12345.txt")
	if err != nil {
		t.Errorf("Read() of non-existent file should not error, got: %v", err)
	}

	if len(lines) != 0 {
		t.Errorf("Read() of non-existent file should return empty slice, got %d lines", len(lines))
	}
}

func TestReadExistingFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test-dotfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write test content
	content := "line1\nline2\nline3\n"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Read the file
	handler := NewHandler()
	lines, err := handler.Read(tmpfile.Name())
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if len(lines) != 3 {
		t.Errorf("Read() got %d lines, want 3", len(lines))
	}

	expectedLines := []string{"line1", "line2", "line3"}
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Read() line %d = %v, want %v", i, line, expectedLines[i])
		}
	}
}

func TestWriteCreatesBackup(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test-dotfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer os.Remove(tmpfile.Name() + ".bak")

	// Write initial content
	initialContent := "initial content\n"
	if _, err := tmpfile.WriteString(initialContent); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Write new content
	handler := NewHandler()
	newLines := []string{"new line 1", "new line 2"}
	if err := handler.Write(tmpfile.Name(), newLines); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Check that backup was created
	backupPath := tmpfile.Name() + ".bak"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("Write() should create backup file")
	}

	// Verify backup contains original content
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(backupContent) != initialContent {
		t.Errorf("Backup content = %v, want %v", string(backupContent), initialContent)
	}

	// Verify new content was written
	newContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	expectedContent := "new line 1\nnew line 2\n"
	if string(newContent) != expectedContent {
		t.Errorf("New content = %v, want %v", string(newContent), expectedContent)
	}
}

func TestBackup(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test-dotfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer os.Remove(tmpfile.Name() + ".bak")

	content := "test content"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Create backup
	handler := NewHandler()
	if err := handler.Backup(tmpfile.Name()); err != nil {
		t.Fatalf("Backup() error = %v", err)
	}

	// Verify backup exists and has same content
	backupPath := tmpfile.Name() + ".bak"
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(backupContent) != content {
		t.Errorf("Backup content = %v, want %v", string(backupContent), content)
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name      string
		shellType shell.ShellType
		wantContains string
	}{
		{
			name:      "create zsh file",
			shellType: shell.Zsh,
			wantContains: "Zsh configuration",
		},
		{
			name:      "create bash file",
			shellType: shell.Bash,
			wantContains: "Bash configuration",
		},
	}

	handler := NewHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary path
			tmpPath := filepath.Join(os.TempDir(), "test-create-"+tt.name)
			defer os.Remove(tmpPath)

			// Create the file
			if err := handler.Create(tmpPath, tt.shellType); err != nil {
				t.Fatalf("Create() error = %v", err)
			}

			// Verify file was created
			if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
				t.Error("Create() should create file")
			}

			// Verify content
			content, err := os.ReadFile(tmpPath)
			if err != nil {
				t.Fatal(err)
			}

			contentStr := string(content)
			if len(contentStr) == 0 {
				t.Error("Create() should write content to file")
			}

			// Check for shell-specific header
			if tt.wantContains != "" {
				if !contains(contentStr, tt.wantContains) {
					t.Errorf("Create() content should contain %v", tt.wantContains)
				}
			}
		})
	}
}

func TestEnsureExists(t *testing.T) {
	handler := NewHandler()
	tmpPath := filepath.Join(os.TempDir(), "test-ensure-exists")
	defer os.Remove(tmpPath)

	// File doesn't exist, should create it
	if err := handler.EnsureExists(tmpPath, shell.Zsh); err != nil {
		t.Fatalf("EnsureExists() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		t.Error("EnsureExists() should create file")
	}

	// Call again, should not error
	if err := handler.EnsureExists(tmpPath, shell.Zsh); err != nil {
		t.Errorf("EnsureExists() on existing file error = %v", err)
	}
}

func TestReadOrCreate(t *testing.T) {
	handler := NewHandler()
	tmpPath := filepath.Join(os.TempDir(), "test-read-or-create")
	defer os.Remove(tmpPath)

	// File doesn't exist, should create and return header
	lines, err := handler.ReadOrCreate(tmpPath, shell.Zsh)
	if err != nil {
		t.Fatalf("ReadOrCreate() error = %v", err)
	}

	if len(lines) == 0 {
		t.Error("ReadOrCreate() should return header lines for new file")
	}

	// File now exists, should read it
	lines2, err := handler.ReadOrCreate(tmpPath, shell.Zsh)
	if err != nil {
		t.Fatalf("ReadOrCreate() second call error = %v", err)
	}

	if len(lines2) == 0 {
		t.Error("ReadOrCreate() should return content for existing file")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
