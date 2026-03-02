package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nickpricks/ft/internal/core"
)

type Config struct {
	NotesDir string `json:"notes_dir"`
}

var configFilePath string

func init() {
	home, err := os.UserHomeDir()
	if err == nil {
		configFilePath = filepath.Join(home, ".fmd.json")
	}
}

// LoadOrInit reads the config file or prompts the user if it doesn't exist.
func LoadOrInit() error {
	if configFilePath == "" {
		// Fallback to local if home dir isn't found
		core.BaseDir = "notes"
		return nil
	}

	data, err := os.ReadFile(configFilePath)
	if err == nil {
		// Config exists!
		var cfg Config
		if err := json.Unmarshal(data, &cfg); err == nil && cfg.NotesDir != "" {
			core.BaseDir = cfg.NotesDir
			return nil
		}
	}

	// If we're running under `go test`, bypass the prompt automatically.
	if strings.HasSuffix(os.Args[0], ".test") || strings.Contains(os.Args[0], "/_go_build_") {
		core.BaseDir = "notes"
		return nil
	}

	// If we get here, the config doesn't exist or is invalid. Prompt the user.
	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, "Documents", "FeatherTrailNotes")

	fmt.Printf("Welcome to FeatherTrailMD!\n")
	fmt.Printf("It looks like this is your first time running the tool.\n")
	fmt.Printf("Where would you like to store your notes? [%s]: ", defaultDir)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	chosenDir := defaultDir
	if input != "" {
		chosenDir = input
	}

	// Save the config
	cfg := Config{NotesDir: chosenDir}
	data, err = json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Make sure the chosen directory exists
	if err := os.MkdirAll(chosenDir, 0755); err != nil {
		return fmt.Errorf("failed to create notes directory: %w", err)
	}

	core.BaseDir = chosenDir
	fmt.Printf("Awesome! Your notes will be saved in: %s\n\n", chosenDir)
	return nil
}
