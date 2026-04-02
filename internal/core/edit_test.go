package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEdit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-edit-tests-")
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create dummy note
	dateFolder := filepath.Join(tmpDir, time.Now().Format(time.DateOnly))
	if err := os.MkdirAll(dateFolder, 0755); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	initialContent := "Initial line.\n"
	err = os.WriteFile(filepath.Join(dateFolder, "01_test.md"), []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy note: %v", err)
	}

	// Test append
	appendedText := "This is an appended line."
	_, err = store.Edit("01", appendedText)
	if err != nil {
		t.Fatalf("Edit failed: %v", err)
	}

	// Verify content after edit
	content, err := store.Read("01")
	if err != nil {
		t.Fatalf("Read after edit failed: %v", err)
	}

	expectedContent := initialContent + "\n" + appendedText + "\n"
	if content != expectedContent {
		t.Errorf("expected %q, got %q", expectedContent, content)
	}

	// Test case: editing non-existent file
	_, err = store.Edit("99", "should fail")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error for invalid ID editing, got %v", err)
	}
}

func TestEdit_ReadOnlyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-edit-readonly-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	dateFolder := filepath.Join(tmpDir, time.Now().Format(time.DateOnly))
	if err := os.MkdirAll(dateFolder, 0755); err != nil {
		t.Fatal(err)
	}

	notePath := filepath.Join(dateFolder, "01_readonly.md")
	if err := os.WriteFile(notePath, []byte("read only note\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(notePath, 0444); err != nil {
		t.Skip("cannot change file permissions on this OS")
	}
	defer os.Chmod(notePath, 0644)

	_, err = store.Edit("01", "should fail")
	if err == nil {
		t.Error("expected error editing read-only file, got nil")
	}
}
