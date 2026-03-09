package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestList(t *testing.T) {
	// Setup a temporary base directory for isolated testing
	tmpDir, err := os.MkdirTemp("", "ft-list-tests-")
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Intercept and restore BaseDir
	originalBaseDir := BaseDir
	BaseDir = tmpDir
	defer func() { BaseDir = originalBaseDir }()

	// Test case: Empty directory
	notes, err := List()
	if err != nil {
		t.Fatalf("List failed on empty dir: %v", err)
	}
	if len(notes) != 0 {
		t.Errorf("expected 0 notes, got %d", len(notes))
	}

	// Setup dummy files for testing
	date1 := "2026-03-01"
	date2 := "2026-03-02"

	if err := os.MkdirAll(filepath.Join(tmpDir, date1), 0755); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, date2), 0755); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, date1, "01_note_a.md"), []byte("A"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, date1, "02_note_b.md"), []byte("B"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, date2, "01_note_c.md"), []byte("C"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	// Create a non-md file that should be ignored
	if err := os.WriteFile(filepath.Join(tmpDir, date2, "02_ignore.txt"), []byte("Ignore"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	// Test case: Verify listing and chronological sorting
	// Expected order:
	// 1. 2026-03-02 01_note_c
	// 2. 2026-03-01 02_note_b
	// 3. 2026-03-01 01_note_a
	notes, err = List()
	if err != nil {
		t.Fatalf("List failed with files: %v", err)
	}

	if len(notes) != 3 {
		t.Fatalf("expected 3 notes, got %d", len(notes))
	}

	// Verify order and extraction logic
	if notes[0].Date != date2 || notes[0].ID != "01" || notes[0].Slug != "note_c" {
		t.Errorf("expected newest note first, got Date=%s ID=%s Slug=%s", notes[0].Date, notes[0].ID, notes[0].Slug)
	}

	if notes[1].Date != date1 || notes[1].ID != "02" || notes[1].Slug != "note_b" {
		t.Errorf("expected second note, got Date=%s ID=%s Slug=%s", notes[1].Date, notes[1].ID, notes[1].Slug)
	}

	if notes[2].Date != date1 || notes[2].ID != "01" || notes[2].Slug != "note_a" {
		t.Errorf("expected oldest note last, got Date=%s ID=%s Slug=%s", notes[2].Date, notes[2].ID, notes[2].Slug)
	}
}

func TestList_MissingBaseDir(t *testing.T) {
	originalBaseDir := BaseDir
	BaseDir = "/path/that/does/not/exist/for/testing"
	defer func() { BaseDir = originalBaseDir }()

	notes, err := List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(notes) != 0 {
		t.Errorf("expected 0 notes, got %d", len(notes))
	}
}
