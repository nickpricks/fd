package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadOrInit_CorruptJSON(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-config-test-")
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origConfigPath := configFilePath
	configFilePath = filepath.Join(tmpDir, ".fmd.json")
	defer func() { configFilePath = origConfigPath }()

	if err := os.WriteFile(configFilePath, []byte("{invalid json"), 0600); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	err = LoadOrInit()
	if err == nil || !strings.Contains(err.Error(), "contains invalid JSON") {
		t.Errorf("expected invalid JSON error, got %v", err)
	}
}

func TestLoadOrInit_EmptyNotesDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-config-test-")
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origConfigPath := configFilePath
	configFilePath = filepath.Join(tmpDir, ".fmd.json")
	defer func() { configFilePath = origConfigPath }()

	if err := os.WriteFile(configFilePath, []byte(`{"notes_dir": ""}`), 0600); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	err = LoadOrInit()
	if err == nil || !strings.Contains(err.Error(), "empty notes_dir field") {
		t.Errorf("expected empty notes_dir error, got %v", err)
	}
}
