package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewNoteStore_ValidDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-store-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatalf("NewNoteStore failed: %v", err)
	}

	abs, _ := filepath.Abs(tmpDir)
	if store.Dir() != abs {
		t.Errorf("expected Dir()=%q, got %q", abs, store.Dir())
	}
}

func TestNewNoteStore_EmptyDir(t *testing.T) {
	_, err := NewNoteStore("")
	if err == nil {
		t.Fatal("expected error for empty base directory, got nil")
	}
	if !strings.Contains(err.Error(), "must not be empty") {
		t.Errorf("expected 'must not be empty' error, got: %v", err)
	}
}

func TestNewNoteStore_RelativeDir(t *testing.T) {
	store, err := NewNoteStore("notes")
	if err != nil {
		t.Fatalf("NewNoteStore failed for relative path: %v", err)
	}

	if !filepath.IsAbs(store.Dir()) {
		t.Errorf("expected absolute path, got %q", store.Dir())
	}
}

func TestNoteStore_AddAndRead(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-store-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Add a note
	path, err := store.Add("hello world")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if !strings.HasSuffix(path, "01_hello_world.md") {
		t.Errorf("unexpected path: %q", path)
	}

	// Read it back
	content, err := store.Read("01")
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if strings.TrimSpace(content) != "hello world" {
		t.Errorf("expected 'hello world', got %q", strings.TrimSpace(content))
	}

	// List
	notes, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
	if notes[0].ID != "01" {
		t.Errorf("expected ID='01', got %q", notes[0].ID)
	}

	// Edit
	editPath, err := store.Edit("01", "appended text")
	if err != nil {
		t.Fatalf("Edit failed: %v", err)
	}
	if editPath != path {
		t.Errorf("expected edit path %q, got %q", path, editPath)
	}

	// Read after edit
	content, err = store.Read("01")
	if err != nil {
		t.Fatalf("Read after edit failed: %v", err)
	}
	if !strings.Contains(content, "appended text") {
		t.Errorf("expected content to contain 'appended text', got %q", content)
	}
}
