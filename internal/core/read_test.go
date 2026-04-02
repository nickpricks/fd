package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-read-tests-")
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create a dummy note to read
	dateFolder := filepath.Join(tmpDir, time.Now().Format(time.DateOnly))
	if err := os.MkdirAll(dateFolder, 0755); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	testContent := "This is a dummy read test\n"
	err = os.WriteFile(filepath.Join(dateFolder, "01_test.md"), []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy note: %v", err)
	}

	// Test successful read
	content, err := store.Read("01")
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if content != testContent {
		t.Errorf("expected %q, got %q", testContent, content)
	}

	// Test case: file not found
	_, err = store.Read("99")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error for invalid ID, got %v", err)
	}
}
