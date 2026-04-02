package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nickpricks/ft/internal/core"
)

func setupConfigTest(t *testing.T) (tmpDir string, cleanup func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "ft-config-test-")
	if err != nil {
		t.Fatal(err)
	}

	origConfigPath := configFilePath
	origHomeDir := homeDir
	origInitErr := initErr
	origBaseDir := core.BaseDir
	origArgs0 := os.Args[0]

	cleanup = func() {
		os.Args[0] = origArgs0
		configFilePath = origConfigPath
		homeDir = origHomeDir
		initErr = origInitErr
		core.BaseDir = origBaseDir
		os.RemoveAll(tmpDir)
	}

	// Override os.Args[0] so isTestRun() returns false, allowing real config logic to run.
	os.Args[0] = "/usr/bin/ft"

	homeDir = tmpDir
	initErr = nil
	configFilePath = filepath.Join(tmpDir, ".fmd.json")

	return tmpDir, cleanup
}

func TestLoadOrInit_ValidConfig(t *testing.T) {
	tmpDir, cleanup := setupConfigTest(t)
	defer cleanup()

	notesDir := filepath.Join(tmpDir, "my-notes")
	cfg := Config{NotesDir: notesDir}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		t.Fatal(err)
	}

	if err := LoadOrInit(); err != nil {
		t.Fatalf("LoadOrInit failed: %v", err)
	}

	if core.BaseDir != notesDir {
		t.Errorf("expected BaseDir=%q, got %q", notesDir, core.BaseDir)
	}
}

func TestLoadOrInit_CorruptJSON(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	if err := os.WriteFile(configFilePath, []byte("{broken json"), 0600); err != nil {
		t.Fatal(err)
	}

	err := LoadOrInit()
	if err == nil {
		t.Fatal("expected error for corrupt JSON, got nil")
	}
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("expected error mentioning 'invalid JSON', got: %v", err)
	}
}

func TestLoadOrInit_EmptyNotesDir(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	data, _ := json.Marshal(Config{NotesDir: ""})
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		t.Fatal(err)
	}

	err := LoadOrInit()
	if err == nil {
		t.Fatal("expected error for empty notes_dir, got nil")
	}
	if !strings.Contains(err.Error(), "empty notes_dir") {
		t.Errorf("expected error mentioning 'empty notes_dir', got: %v", err)
	}
}

func TestLoadOrInit_InitError(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	initErr = errors.New("simulated home dir failure")

	err := LoadOrInit()
	if err == nil {
		t.Fatal("expected error when initErr is set, got nil")
	}
	if !strings.Contains(err.Error(), "cannot load config") {
		t.Errorf("expected 'cannot load config' error, got: %v", err)
	}
}

func TestLoadOrInit_ConfigRoundtrip(t *testing.T) {
	tmpDir, cleanup := setupConfigTest(t)
	defer cleanup()

	notesDir := filepath.Join(tmpDir, "roundtrip-notes")
	cfg := Config{NotesDir: notesDir}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		t.Fatal(err)
	}

	if err := LoadOrInit(); err != nil {
		t.Fatalf("first LoadOrInit failed: %v", err)
	}
	if core.BaseDir != notesDir {
		t.Errorf("first load: expected %q, got %q", notesDir, core.BaseDir)
	}

	core.BaseDir = "notes"
	if err := LoadOrInit(); err != nil {
		t.Fatalf("second LoadOrInit failed: %v", err)
	}
	if core.BaseDir != notesDir {
		t.Errorf("second load: expected %q, got %q", notesDir, core.BaseDir)
	}
}

func TestLoadOrInit_PermissionDenied(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	data, _ := json.Marshal(Config{NotesDir: "/tmp/notes"})
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(configFilePath, 0000); err != nil {
		t.Skip("cannot change file permissions on this OS")
	}
	defer os.Chmod(configFilePath, 0600)

	err := LoadOrInit()
	if err == nil {
		t.Fatal("expected error for unreadable config, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("expected 'failed to read config file' error, got: %v", err)
	}
}
