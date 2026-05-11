package config

import (
	"os"
	"path/filepath"
	"testing"
)

func withTempConfigDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("AM_CONFIG_DIR", dir)
	return dir
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	withTempConfigDir(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.Sync != nil {
		t.Fatalf("expected empty Sync, got %+v", cfg.Sync)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	dir := withTempConfigDir(t)

	in := &Config{Sync: &SyncConfig{
		Provider:     "gist",
		GistID:       "abc123",
		GistFilename: "am-aliases.json",
		LastSynced:   "2025-01-01T00:00:00Z",
		LastHash:     "deadbeef",
	}}
	if err := Save(in); err != nil {
		t.Fatalf("Save: %v", err)
	}

	path := filepath.Join(dir, "config.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file at %s: %v", path, err)
	}

	out, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if out.Sync == nil {
		t.Fatal("expected Sync to be loaded")
	}
	if *out.Sync != *in.Sync {
		t.Fatalf("round-trip mismatch:\n want: %+v\n got:  %+v", in.Sync, out.Sync)
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "nested", "am")
	t.Setenv("AM_CONFIG_DIR", target)

	if err := Save(&Config{Sync: &SyncConfig{Provider: "gist"}}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, "config.json")); err != nil {
		t.Fatalf("expected file created: %v", err)
	}
}

func TestDirHonorsXDG(t *testing.T) {
	os.Unsetenv("AM_CONFIG_DIR")
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	dir, err := Dir()
	if err != nil {
		t.Fatalf("Dir: %v", err)
	}
	want := filepath.Join(xdg, "am")
	if dir != want {
		t.Fatalf("Dir = %q, want %q", dir, want)
	}
}
