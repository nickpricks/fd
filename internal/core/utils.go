// Package notes provides the core filesystem logic for FeatherTrailMD.
package core

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nickpricks/ft/internal/constants"
)

// GetDateFolder returns the absolute or relative path to today's folder.
// e.g., "notes/2026-03-01"
func GetDateFolder() string {
	return filepath.Join(BaseDir, time.Now().Format(time.DateOnly))
}

// Slugify converts an arbitrary string of text into a filesystem-safe string.
// It lowercases everything, removes special characters, and takes the first 5 words.
func Slugify(text string) string {
	text = strings.ToLower(text)
	reg := regexp.MustCompile("[^a-z0-9]+")
	text = reg.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	words := strings.Split(text, " ")

	limit := 5
	if len(words) < limit {
		limit = len(words)
	}

	// Fallback if the string was completely empty or stripped of all valid chars
	if limit == 0 || words[0] == "" {
		return "note"
	}
	return strings.Join(words[:limit], "_")
}

// findNoteByID recursively checks today's folder, then previous days' folders,
// to locate a file whose name starts with the given ID prefix.
func findNoteByID(id string) (string, error) {
	if _, err := os.Stat(BaseDir); os.IsNotExist(err) {
		return "", fmt.Errorf(constants.ErrNotesDirNotFound)
	}

	dateFolders, err := os.ReadDir(BaseDir)
	if err != nil {
		return "", err
	}

	// Sort folders descending to prioritize reading newer notes
	sort.Slice(dateFolders, func(i, j int) bool {
		return dateFolders[i].Name() > dateFolders[j].Name()
	})

	for _, folder := range dateFolders {
		if !folder.IsDir() {
			continue
		}

		folderPath := filepath.Join(BaseDir, folder.Name())
		files, err := os.ReadDir(folderPath)
		if err != nil {
			continue // Skip unreadable folders
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasPrefix(file.Name(), id+"_") {
				return filepath.Join(folderPath, file.Name()), nil
			}
		}
	}
	return "", fmt.Errorf(constants.ErrNoteNotFound, id)
}

// getNextID scans the given date folder to find the highest existing ID prefix,
// then returns the next available ID formatted as a two-digit string (e.g. "01", "02").
func getNextID(dateFolder string) (string, error) {
	entries, err := os.ReadDir(dateFolder)
	if err != nil {
		// If the directory doesn't exist yet, it's the first note.
		if os.IsNotExist(err) {
			return "01", nil
		}
		return "", err
	}

	maxID := 0
	for _, entry := range entries {
		// Only look at markdown files
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), constants.ExtMD) {
			continue
		}

		// Split by underscore to extract the ID prefix
		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) > 0 {
			if id, err := strconv.Atoi(parts[0]); err == nil {
				if id > maxID {
					maxID = id
				}
			}
		}
	}
	return fmt.Sprintf("%02d", maxID+1), nil
}
