package core

import (
	"os"
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

func TestGetNextID(t *testing.T) {
	// Setup a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "notes-test-")
	if err != nil {
		t.Fatal(err)
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

	// Create some dummy markdown files
	os.WriteFile(tmpDir+"/01_test.md", []byte("test"), 0644)
	os.WriteFile(tmpDir+"/02_test.md", []byte("test"), 0644)
	os.WriteFile(tmpDir+"/not_a_note.txt", []byte("test"), 0644)

	// Test with existing files
	id, err = getNextID(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "03" {
		t.Errorf("expected 03, got %v", id)
	}
}
