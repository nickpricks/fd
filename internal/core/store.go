// Package notes provides the core filesystem logic for FeatherTrailMD.
package core

import (
	"errors"
	"fmt"
	"path/filepath"
)

// NoteStore encapsulates note storage operations rooted at a specific directory.
type NoteStore struct {
	baseDir string
}

// NewNoteStore creates a NoteStore after validating the base directory.
func NewNoteStore(baseDir string) (*NoteStore, error) {
	if baseDir == "" {
		return nil, errors.New("base directory must not be empty")
	}
	abs, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("invalid base directory: %w", err)
	}
	return &NoteStore{baseDir: abs}, nil
}

// Dir returns the resolved base directory path.
func (s *NoteStore) Dir() string {
	return s.baseDir
}

// Add creates a new note in today's date folder.
func (s *NoteStore) Add(text string) (string, error) {
	return addNote(s.baseDir, text)
}

// List returns all notes chronologically sorted.
func (s *NoteStore) List() ([]NoteInfo, error) {
	return listNotes(s.baseDir)
}

// Read returns the raw content of a note by ID.
func (s *NoteStore) Read(id string) (string, error) {
	return readNote(s.baseDir, id)
}

// Edit appends text to an existing note by ID.
func (s *NoteStore) Edit(id string, text string) (string, error) {
	return editNote(s.baseDir, id, text)
}
