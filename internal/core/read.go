// Package notes provides the core filesystem logic for FeatherTrailMD.
package core

import (
	"os"
)

// Read locates a note by its ID and returns its raw string content.
// It searches in the current day's folder first, then works backwards.
func Read(id string) (string, error) {
	path, err := findNoteByID(id)
	if err != nil {
		return "", err
	}

	// Read the entire file content into memory
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
