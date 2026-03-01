// Package notes provides the core filesystem logic for FeatherTrailMD.
package core

import (
	"os"

	"github.com/nickpricks/ft/internal/constants"
)

// Edit locates a note by its ID and appends the provided text to the bottom of the file.
// Editing in this version is strictly limited to an append operation.
func Edit(id string, text string) (string, error) {
	path, err := findNoteByID(id)
	if err != nil {
		return "", err
	}

	// Open file in append mode. Create it if it doesn't exist (though findNoteByID ensures it does).
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, constants.FilePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Append text with a leading newline for separation
	if _, err := f.WriteString("\n" + text + "\n"); err != nil {
		return "", err
	}

	return path, nil
}
