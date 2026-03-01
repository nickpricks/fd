package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEdit(t *testing.T) {
	// Setup isolated tmp dir
	tmpDir, err := os.MkdirTemp("", "ft-edit-tests-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	originalBaseDir := BaseDir
	BaseDir = tmpDir
	defer func() { BaseDir = originalBaseDir }()

	// Create dummy note
	dateFolder := GetDateFolder()
	os.MkdirAll(dateFolder, 0755)

	initialContent := "Initial line.\n"
	err = os.WriteFile(filepath.Join(dateFolder, "01_test.md"), []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy note: %v", err)
	}

	// Test append
	appendedText := "This is an appended line."
	_, err = Edit("01", appendedText)
	if err != nil {
		t.Fatalf("Edit failed: %v", err)
	}

	// Verify content after edit
	content, err := Read("01")
	if err != nil {
		t.Fatalf("Read after edit failed: %v", err)
	}

	expectedContent := initialContent + "\n" + appendedText + "\n"
	if content != expectedContent {
		t.Errorf("expected %q, got %q", expectedContent, content)
	}

	// Test case: editing non-existent file
	_, err = Edit("99", "should fail")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error for invalid ID editing, got %v", err)
	}
}
