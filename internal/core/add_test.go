package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-tests-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	content := "This is a test note for this version"

	// Test Add
	path, err := store.Add(content)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if !strings.HasSuffix(path, "01_this_is_a_test_note.md") {
		t.Errorf("expected path to end with '01_this_is_a_test_note.md', got %q", path)
	}

	// Verify file content
	readBack, err := store.Read("01")
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if strings.TrimSpace(readBack) != content {
		t.Errorf("expected content %q, got %q", content, strings.TrimSpace(readBack))
	}
}

func TestAdd_MultipleNotes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-tests-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	content1 := "First note"
	content2 := "Second note"

	path1, err := store.Add(content1)
	if err != nil {
		t.Fatalf("Add first note failed: %v", err)
	}
	if !strings.HasSuffix(path1, "01_first_note.md") {
		t.Errorf("expected path to end with '01_first_note.md', got %q", path1)
	}

	path2, err := store.Add(content2)
	if err != nil {
		t.Fatalf("Add second note failed: %v", err)
	}
	if !strings.HasSuffix(path2, "02_second_note.md") {
		t.Errorf("expected path to end with '02_second_note.md', got %q", path2)
	}
}

func TestAdd_EmptyInput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-tests-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	path, err := store.Add("")
	if err != nil {
		t.Fatalf("Add empty note failed: %v", err)
	}

	if !strings.HasSuffix(path, "01_note.md") {
		t.Errorf("expected path to end with '01_note.md', got %q", path)
	}
}

func TestAdd_Fail(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-tests-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file where the baseDir should be, so MkdirAll fails
	blockingFile := filepath.Join(tmpDir, "blocked_dir")
	if err := os.WriteFile(blockingFile, []byte("blocker"), 0644); err != nil {
		t.Fatal(err)
	}

	store, err := NewNoteStore(blockingFile)
	if err != nil {
		t.Fatal(err)
	}

	// This should fail because it attempts to create a directory inside a file
	_, err = store.Add("This should fail")
	if err == nil {
		t.Errorf("expected Add to fail when baseDir is a file, but it succeeded")
	}
}
