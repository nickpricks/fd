package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	// Setup isolated tmp dir
	tmpDir, err := os.MkdirTemp("", "ft-read-tests-")
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalBaseDir := BaseDir
	BaseDir = tmpDir
	defer func() { BaseDir = originalBaseDir }()

	// Create a dummy note to read
	dateFolder := GetDateFolder()
	if err := os.MkdirAll(dateFolder, 0755); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	testContent := "This is a dummy read test\n"
	err = os.WriteFile(filepath.Join(dateFolder, "01_test.md"), []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy note: %v", err)
	}

	// Test successful read
	content, err := Read("01")
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if content != testContent {
		t.Errorf("expected %q, got %q", testContent, content)
	}

	// Test case: file not found
	_, err = Read("99")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error for invalid ID, got %v", err)
	}
}
