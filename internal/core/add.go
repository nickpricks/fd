// Package notes provides the core filesystem logic for FeatherTrailMD.
package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nickpricks/ft/internal/constants"
)

// Add creates a new note in today's date folder with the provided text.
// It returns the path to the newly created file or an error if the operation fails.
func Add(text string) (string, error) {
	dateFolder := GetDateFolder()

	// Ensure the directory exists
	if err := os.MkdirAll(dateFolder, constants.DirPerm); err != nil {
		return "", fmt.Errorf(constants.ErrCreateDateFolder, err)
	}

	id, err := getNextID(dateFolder)
	if err != nil {
		return "", fmt.Errorf(constants.ErrGenerateID, err)
	}

	slug := Slugify(text)
	filename := fmt.Sprintf("%s_%s%s", id, slug, constants.ExtMD)
	path := filepath.Join(dateFolder, filename)

	// Create the file and write the text
	if err := os.WriteFile(path, []byte(text+"\n"), constants.FilePerm); err != nil {
		return "", fmt.Errorf(constants.ErrWriteNote, err)
	}
	return path, nil
}
