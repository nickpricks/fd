// Package notes provides the core filesystem logic for FeatherTrailMD.
package core

import (
	"fmt"
	"os"
)

// Read locates a note by its ID and returns its raw string content.
func Read(id string) (string, error) {
	path, err := findNoteByID(id)
	if err != nil {
		return "", fmt.Errorf("failed to find note: %w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read note content: %w", err)
	}

	return string(content), nil
}
