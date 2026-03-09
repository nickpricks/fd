package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello_world"},
		{"This is a very long sentence that should be cut", "this_is_a_very_long"},
		{"Special!@# Characters???", "special_characters"},
		{"", "note"},
		{"---", "note"},
	}

	for _, tt := range tests {
		result := Slugify(tt.input)
		if result != tt.expected {
			t.Errorf("Slugify(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestFindNoteByID_MissingBaseDir(t *testing.T) {
	originalBaseDir := BaseDir
	BaseDir = "/does/not/exist/surely/ft"
	defer func() { BaseDir = originalBaseDir }()

	_, err := findNoteByID("99")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected notes dir not found error, got %v", err)
	}
}

func TestGetNextID(t *testing.T) {
	// Setup a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "notes-test-")
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test empty directory
	id, err := getNextID(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "01" {
		t.Errorf("expected 01 for empty dir, got %v", id)
	}

	if err := os.WriteFile(tmpDir+"/01_test.md", []byte("test"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.WriteFile(tmpDir+"/02_test.md", []byte("test"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.WriteFile(tmpDir+"/not_a_note.txt", []byte("test"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	// Test with existing files
	id, err = getNextID(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "03" {
		t.Errorf("expected 03, got %v", id)
	}
}

func TestGetNextID_NonexistentDir(t *testing.T) {
	// Setup a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "notes-test-")
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Subdir that does NOT exist
	subDir := filepath.Join(tmpDir, "missing")
	id, err := getNextID(subDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "01" {
		t.Errorf("expected 01 for missing dir, got %v", id)
	}
}
