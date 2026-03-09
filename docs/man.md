# FeatherTrailMD (`ft`) Developer Manual

This document provides a line-by-line breakdown of the core CLI execution flow to help new contributors understand how `ft` works under the hood.

## The CLI Framework: Cobra
We use [spf13/cobra](https://github.com/spf13/cobra), a popular Go library for creating powerful modern CLI applications. 
- A `cobra.Command` struct defines everything about a command: its name (`Use`), description (`Short`/`Long`), examples (`Example`), arguments validation (`Args`), and the actual logic to execute (`RunE`).
- Commands are arranged in a tree. The root command is `ft`, and subcommands like `add` or `list` are added via `rootCmd.AddCommand(addCmd)`.

---

## `cmd/feathertrailmd/main.go`
*(Note: The CLI entry point was renamed from `cmd/ft` to `cmd/feathertrailmd` because the `.gitignore` rule for the compiled binary `ft` was accidentally ignoring the entire source folder, breaking the GitHub Actions automated release build.)*

This is the entry point of the compiled binary.

```go
package main

import (
	"fmt"
	"os"
	"github.com/nickpricks/ft/internal/cli"
)

func main() {
    // 1. We call cli.Execute(), which delegates to our root Cobra command.
	if err := cli.Execute(); err != nil {
        // 2. If Execute returns an error (e.g., bad arguments or underlying failure), 
        // we print it to the standard error stream.
		fmt.Fprintln(os.Stderr, err)
        // 3. We exit with a non-zero status code to indicate the process failed.
		os.Exit(1)
	}
}
```

---

## `internal/cli/root.go`
This defines the base `ft` command.

```go
package cli

import (
    "github.com/nickpricks/ft/internal/config"
    "github.com/spf13/cobra"
)

// 1. Define the root command struct.
var rootCmd = &cobra.Command{
	Use:   "ft", // The actual command users type
	Short: "FeatherTrailMD is a quick notes tool - Super Fast Thoughts Notes", // Brief description
	Long:  `FeatherTrailMD (ft) is a simple, filesystem-first notes assistant.`, // Detailed description
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" || cmd.CalledAs() == "help" {
			return nil
		}
		return config.LoadOrInit()
	},
}

// 2. Execute exposes the private rootCmd.Execute() to other packages (like our main.go).
func Execute() error {
	return rootCmd.Execute()
}

// 3. init() runs automatically before main(). This is where we wire up global flags.
func init() {
	// e.g., rootCmd.PersistentFlags().String("config", "", "config file...")
}
```

---

## `internal/cli/add.go`
This handles the `ft add` command.

```go
package cli

import (
	"fmt"
	"strings"
	"github.com/nickpricks/ft/internal/core"
	"github.com/spf13/cobra"
)

// 1. Define the 'add' subcommand.
var addCmd = &cobra.Command{
	Use:     "add [text...]", // The "..." indicates it accepts multiple words
	Short:   "Quickly capture a new note",
	Long:    `The add command takes any text you provide and creates a new markdown note...`,
	Example: `  ft add "Meeting notes"`,
	Args:    cobra.MinimumNArgs(1), // Validation: Fails if the user provides 0 arguments.
	
	// 2. RunE is the actual execution function that returns an error.
	RunE: func(cmd *cobra.Command, args []string) error {
        // 3. args is a slice of strings. If the user didn't use quotes (ft add Hello World), 
        // args will be ["Hello", "World"]. We join them back together with spaces.
		text := strings.Join(args, " ")
        
        // 4. Call our core business logic in the 'notes' package.
		path, err := notes.Add(text)
		if err != nil {
			return err // Return error to Cobra, which will print it.
		}
        
        // 5. Success! Print the path to the newly created note.
		fmt.Printf("Note created: %s\n", path)
		return nil
	},
}

// 6. Automatically register addCmd under rootCmd when the program starts.
func init() {
	rootCmd.AddCommand(addCmd)
}
```
*(Similar Cobra definitions exist for `list.go`, `read.go`, and `edit.go`, each parsing their specific arguments and calling the corresponding `notes` package function).*

---

## `internal/core/add.go`
This is the core business logic where notes are actually created on the filesystem.

```go
package core

import (
	"fmt"
	"os"
	"path/filepath"
	// ... other imports
)

// Add creates a new note in today's date folder with the provided text.
func Add(text string) (string, error) {
    // 1. Get today's folder path (e.g., notes/2026-03-01)
	dateFolder := GetDateFolder()

	// 2. Ensure the directory exists, creating it if necessary
	if err := os.MkdirAll(dateFolder, 0755); err != nil {
		return "", fmt.Errorf("failed to create date folder: %w", err)
	}

    // 3. Scan the folder to find the next available ID (e.g. "01", "02")
	id, err := getNextID(dateFolder)
	if err != nil {
        return "", err
    }

    // 4. Generate a clean URL-friendly slug from the text
	slug := Slugify(text)
	filename := fmt.Sprintf("%s_%s.md", id, slug)
	path := filepath.Join(dateFolder, filename)

	// 5. Create the file and write the text
	if err := os.WriteFile(path, []byte(text+"\n"), 0644); err != nil {
		return "", fmt.Errorf("failed to write note: %w", err)
	}
    
    // 6. Return the path so the CLI can print it
	return path, nil
}
```
---

## `internal/core/list.go`
This handles retrieving and sorting all notes.
- Uses `filepath.WalkDir` to efficiently scan the base directory.
- Filters out non-markdown files and directories.
- Parses the date, ID, and slug from filename structures (`YYYY-MM-DD/ID_slug.md`).
- Sorts notes descending by date, then descending by ID, so the newest notes appear first.

---

## `internal/core/read.go`
This fetches a note's raw text given its ID.
- Reuses `findNoteByID(id)` to scan the latest day folders first.
- Uses `os.ReadFile` to pull the raw string content into memory.

---

## `internal/core/edit.go`
This appends text to an existing note.
- Focuses solely on append operations for safety.
- Opens the specific note file with `os.O_APPEND|os.O_WRONLY` flags.
- Appends the new text securely to the end of the file.

---

## `internal/core/add_test.go`
This is how we test our note-creation logic. We isolate our tests from the real filesystem so that running tests doesn't clutter the user's actual notes.

```go
package core

import (
	"os"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	// 1. Create a temporary isolated directory for the test to use.
	tmpDir, err := os.MkdirTemp("", "ft-tests-")
	if err != nil {
		t.Fatal(err)
	}
	// 2. Schedule the temporary directory for deletion when the test finishes.
	defer os.RemoveAll(tmpDir)

	// 3. Override the global BaseDir to point to our isolated tmpDir.
	// We save the original BaseDir and use a deferred function to restore it later.
	originalBaseDir := BaseDir
	BaseDir = tmpDir
	defer func() { BaseDir = originalBaseDir }()

    // 4. Test the business logic
	content := "This is a test note for this version"
	path, err := Add(content)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

    // 5. Verify the generated filename is correct (01_ + slugified content)
	if !strings.HasSuffix(path, "01_this_is_a_test_note.md") {
		t.Errorf("expected path to end with '01_this_is_a_test_note.md', got %q", path)
	}
    
    // ... Additional test logic to verify the file content was written correctly ...
}
```

*(Other test files like `list_test.go`, `read_test.go`, `edit_test.go`, and `utils_test.go` follow this exact same pattern of using `os.MkdirTemp` and overriding `BaseDir` to safely isolate tests from the real notes directory)*

---

## `internal/constants/constants.go`
To prevent magic strings and hardcoded permission bits from scattering across the codebase, we centralize all shared values here:

```go
package constants

import "os"

// File system constants
const (
	FilePerm os.FileMode = 0644
	DirPerm  os.FileMode = 0755
	ExtMD                = ".md"
)

// Log messages
const (
	LogNoteCreated = "Note created: %s\n"
	LogNoNotes     = "No notes found."
	// ...
)
```

By referencing `constants.FilePerm` rather than `0644`, the codebase remains clean, consistent, and easily updatable.

---

## `internal/core/utils.go`
This file acts as a shared utility belt for the notes package. Instead of duplicating logic across `add.go` and `list.go`, we define reusable functions here:

- **`GetDateFolder()`**: Uses `time.Now()` to dynamically generate today's folder path (e.g. `notes/2026-03-01`).
- **`Slugify(text)`**: Safely converts the first few words of a note into a valid file slug by stripping punctuation and injecting underscores.
- **`getNextID(folder)`**: Scans a specific date folder and returns the next sequential ID, ensuring `01`, `02`, `03` incrementation logic lives in exactly one place.
