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
var homeDir string
var initErr error

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		initErr = fmt.Errorf("could not determine home directory: %w", err)
		return
	}
	homeDir = home
	configFilePath = filepath.Join(homeDir, ".fmd.json")
}

// LoadOrInit reads the config file or prompts the user if it doesn't exist.
func LoadOrInit() error {
	if initErr != nil {
		return fmt.Errorf("cannot load config: %w", initErr)
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read config file %s: %w", configFilePath, err)
		}
		// File genuinely doesn't exist -> first run, fall through to prompt
	} else {
		var cfg Config
		if jsonErr := json.Unmarshal(data, &cfg); jsonErr != nil {
			return fmt.Errorf("config file %s contains invalid JSON: %w\nTo fix: delete the file and re-run ft", configFilePath, jsonErr)
		}
		if cfg.NotesDir == "" {
			return fmt.Errorf("config file %s has empty notes_dir field", configFilePath)
		}
		core.BaseDir = cfg.NotesDir
		return nil
	}

	// If we're running under `go test`, bypass the prompt automatically.
	// Note: The forward slash check may not match on Windows.
	exe := os.Args[0]
	base := filepath.Base(exe)
	if strings.HasSuffix(base, ".test") || strings.HasSuffix(base, ".test.exe") ||
		strings.Contains(exe, string(filepath.Separator)+"_go_build_") {
		core.BaseDir = "notes"
		return nil
	}

	// If we get here, the config file is missing. Prompt the user.
	defaultDir := filepath.Join(homeDir, "Documents", "FeatherTrailNotes")

	fmt.Printf("Welcome to FeatherTrailMD!\n")
	fmt.Printf("It looks like this is your first time running the tool.\n")
	fmt.Printf("Where would you like to store your notes? [%s]: ", defaultDir)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input (is stdin available?): %w", err)
	}
	input = strings.TrimSpace(input)

	chosenDir := defaultDir
	if input != "" {
		chosenDir = input
	}

	chosenDir, err = filepath.Abs(chosenDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Save the config
	cfg := Config{NotesDir: chosenDir}
	data, err = json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
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
