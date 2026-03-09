// Package notes provides the core filesystem logic for FeatherTrailMD.
package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nickpricks/ft/internal/constants"
)

// List recursively walks the BaseDir and returns a chronological slice of all NoteInfo.
// If the BaseDir does not exist, it safely returns an empty slice.
func List() ([]NoteInfo, error) {
	var notes []NoteInfo

	// If base directory doesn't exist, there are simply no notes yet.
	if _, err := os.Stat(BaseDir); os.IsNotExist(err) {
		return notes, nil
	}

	err := filepath.WalkDir(BaseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if d.IsDir() || !strings.HasSuffix(d.Name(), constants.ExtMD) {
			return nil
		}

		rel, err := filepath.Rel(BaseDir, path)
		if err != nil {
			return fmt.Errorf("failed to compute relative path for %s: %w", path, err)
		}
		parts := strings.Split(filepath.ToSlash(rel), "/") // Expected: YYYY-MM-DD/ID_slug.md

		if len(parts) >= 2 {
			date := parts[len(parts)-2]
			filename := parts[len(parts)-1]

			// Extract ID and slug from filename
			fileParts := strings.SplitN(filename, "_", 2)
			id := fileParts[0]
			slug := "note"
			if len(fileParts) > 1 {
				slug = strings.TrimSuffix(fileParts[1], constants.ExtMD)
			}

			notes = append(notes, NoteInfo{
				Path: path,
				Date: date,
				ID:   id,
				Slug: slug,
			})
		}
		return nil
	})

	// Sort by date (descending), then ID (descending) so newest appears first
	sort.Slice(notes, func(i, j int) bool {
		if notes[i].Date != notes[j].Date {
			return notes[i].Date > notes[j].Date
		}
		return notes[i].ID > notes[j].ID
	})

	return notes, err
}
