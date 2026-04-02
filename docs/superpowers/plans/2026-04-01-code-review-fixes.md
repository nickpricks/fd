# Code Review Board Clearance — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix all 25 items from `docs/CODE_REVIEW.md` — critical config bugs, error handling gaps, test coverage, documentation, and architecture improvements.

**Architecture:** Work progresses foundation-up: sentinel errors first, then config hardening, core error handling, tests, docs, and finally the NoteStore structural refactor. Each task produces a commit. TDD where applicable.

**Tech Stack:** Go 1.26, Cobra CLI, stdlib only (no new dependencies)

**Items already addressed:** D2 (man.md code listings) and D3 (ref.md config section) were fixed in a prior docs update pass.

---

## File Map

| File | Action | Tasks |
|------|--------|-------|
| `internal/core/errors.go` | Create | 1 |
| `internal/constants/constants.go` | Modify | 1 |
| `internal/config/config.go` | Modify | 2, 3 |
| `internal/cli/root.go` | Modify | 3, 9 |
| `internal/core/utils.go` | Modify | 4 |
| `internal/core/list.go` | Modify | 4 |
| `internal/core/edit.go` | Modify | 5 |
| `internal/core/read.go` | Modify | 5 |
| `internal/cli/list.go` | Modify | 5 |
| `internal/config/config_test.go` | Create | 6 |
| `internal/core/utils_test.go` | Modify | 7 |
| `internal/core/edit_test.go` | Modify | 7 |
| `internal/core/list_test.go` | Modify | 7 |
| `internal/core/read_test.go` | Modify | 7 |
| `internal/core/types.go` | Modify | 8 |
| `README.md` | Modify | 8 |
| `cmd/feathertrailmd/README.md` | Modify | 8 |
| `internal/core/store.go` | Create | 9 |
| `internal/core/store_test.go` | Create | 9 |
| `internal/core/add.go` | Modify | 9 |
| `internal/core/add_test.go` | Modify | 9 |
| All `*_test.go` in core | Modify | 9 |
| All `internal/cli/*.go` | Modify | 9 |
| `internal/core/noteid.go` | Create | 10 |
| `internal/core/noteid_test.go` | Create | 10 |

---

### Task 1: Add Sentinel Errors (A4)

**Covers:** A4
**Files:**
- Create: `internal/core/errors.go`
- Modify: `internal/constants/constants.go` (remove string error templates)
- Modify: `internal/core/utils.go:78` (use sentinel)
- Modify: `internal/core/add.go:18-19` (use sentinel)

- [ ] **Step 1: Create `internal/core/errors.go` with sentinel errors**

```go
package core

import "errors"

var (
	ErrNoteNotFound    = errors.New("note not found")
	ErrNotesDirMissing = errors.New("notes directory not found")
	ErrCreateDateDir   = errors.New("failed to create date folder")
	ErrGenerateID      = errors.New("failed to generate next ID")
	ErrWriteNote       = errors.New("failed to write note")
	ErrMaxNotesPerDay  = errors.New("maximum notes per day (99) exceeded")
)
```

- [ ] **Step 2: Update `findNoteByID` in `utils.go` to use sentinel error**

Change line 78 from:
```go
return "", fmt.Errorf(constants.ErrNoteNotFound, id)
```
To:
```go
return "", fmt.Errorf("%w: %s", ErrNoteNotFound, id)
```

And line 49 from:
```go
return "", fmt.Errorf(constants.ErrNotesDirNotFound)
```
To:
```go
return "", ErrNotesDirMissing
```

- [ ] **Step 3: Update `Add` in `add.go` to use sentinel errors**

Change line 19 from:
```go
return "", fmt.Errorf(constants.ErrCreateDateFolder, err)
```
To:
```go
return "", fmt.Errorf("%w: %w", ErrCreateDateDir, err)
```

Change line 23 from:
```go
return "", fmt.Errorf(constants.ErrGenerateID, err)
```
To:
```go
return "", fmt.Errorf("%w: %w", ErrGenerateID, err)
```

Change line 32 from:
```go
return "", fmt.Errorf(constants.ErrWriteNote, err)
```
To:
```go
return "", fmt.Errorf("%w: %w", ErrWriteNote, err)
```

- [ ] **Step 4: Remove unused error templates from `constants.go`**

Remove the `ErrNotesDirNotFound`, `ErrNoteNotFound`, `ErrCreateDateFolder`, `ErrGenerateID`, `ErrWriteNote` constants from `internal/constants/constants.go` since they're replaced by sentinel errors in `core`.

- [ ] **Step 5: Run tests to verify nothing broke**

Run: `go test ./... -v`
Expected: All existing tests pass. The `strings.Contains(err.Error(), "not found")` checks in `edit_test.go:52` and `read_test.go:44` still pass because the sentinel error message contains "not found".

- [ ] **Step 6: Run vet**

Run: `go vet ./...`
Expected: Clean

- [ ] **Step 7: Commit**

```
feat: add sentinel errors replacing string error templates (A4)
```

---

### Task 2: Config Error Handling (C1-C6)

**Covers:** C1, C2, C3, C4, C5, C6
**Files:**
- Modify: `internal/config/config.go`

- [ ] **Step 1: Rewrite `config.go` — store homeDir, fix all silent fallbacks**

Replace the entire `internal/config/config.go` with:

```go
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

var (
	configFilePath string
	homeDir        string
	initErr        error
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		initErr = fmt.Errorf("could not determine home directory: %w", err)
		return
	}
	homeDir = home
	configFilePath = filepath.Join(home, ".fmd.json")
}

// isTestRun checks if the binary was launched by `go test`.
func isTestRun() bool {
	exe := os.Args[0]
	base := filepath.Base(exe)
	if strings.HasSuffix(base, ".test") || strings.HasSuffix(base, ".test.exe") {
		return true
	}
	if strings.Contains(exe, string(filepath.Separator)+"_go_build_") {
		return true
	}
	return false
}

// LoadOrInit reads the config file or prompts the user if it doesn't exist.
func LoadOrInit() error {
	if initErr != nil {
		return fmt.Errorf("cannot load config: %w", initErr)
	}

	// If running under `go test`, bypass the prompt.
	if isTestRun() {
		core.BaseDir = "notes"
		return nil
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read config file %s: %w", configFilePath, err)
		}
		// File genuinely doesn't exist — first run, fall through to prompt
	} else {
		// File exists — parse it
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

	// First run — prompt user for notes directory
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

	// Resolve to absolute path
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
```

Key changes:
- C1: `ReadFile` errors distinguished from file-not-found via `os.IsNotExist`
- C2: JSON parse errors return explicit error instead of silent re-prompt
- C3: `initErr` stored and surfaced — no silent fallback to relative path
- C4: `reader.ReadString` error checked — fails in non-interactive envs
- C5: `homeDir` stored in package var during `init()`, reused for default dir
- C6: `isTestRun()` checks `.test.exe` and uses `filepath.Separator`
- I6: Config file written with `0600` instead of `0644`
- I7: `filepath.Abs(chosenDir)` called before saving

- [ ] **Step 2: Run tests**

Run: `go test ./... -v`
Expected: All tests pass. Config's test detection still works (isTestRun uses filepath.Separator).

- [ ] **Step 3: Run vet**

Run: `go vet ./...`
Expected: Clean

- [ ] **Step 4: Commit**

```
fix: config error handling — surface errors instead of silent fallbacks (C1-C6, I6, I7)
```

---

### Task 3: Skip Config for Help/Version (I8)

**Covers:** I8
**Files:**
- Modify: `internal/cli/root.go:15-17`

- [ ] **Step 1: Update PersistentPreRunE to skip config for help**

Replace `internal/cli/root.go` content with:

```go
package cli

import (
	"github.com/nickpricks/ft/internal/config"
	"github.com/nickpricks/ft/internal/constants"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     constants.RootUse,
	Short:   constants.RootShort,
	Long:    constants.RootLong,
	Example: constants.RootExample,
	Version: constants.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" || cmd.CalledAs() == "help" {
			return nil
		}
		return config.LoadOrInit()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Root flags can be added here
}
```

- [ ] **Step 2: Verify `ft --help` doesn't trigger config**

Run: `go build -o ft ./cmd/feathertrailmd && ./ft --help`
Expected: Help text displayed without config prompt. Note: Cobra handles `--version` via the `Version` field and doesn't trigger `PersistentPreRunE`.

- [ ] **Step 3: Run tests**

Run: `go test ./... -v`
Expected: All pass

- [ ] **Step 4: Commit**

```
fix: skip config init for help command (I8)
```

---

### Task 4: Core Error Handling — utils.go & list.go (C7, I1, I2, I4)

**Covers:** C7, I1, I2, I4
**Files:**
- Modify: `internal/core/utils.go:47-49, 67-70, 83-111`
- Modify: `internal/core/list.go:19-20, 33`

- [ ] **Step 1: Fix `getNextID` to cap at 99 notes/day (C7)**

In `internal/core/utils.go`, replace lines 109-111:
```go
	return fmt.Sprintf("%02d", maxID+1), nil
```
With:
```go
	next := maxID + 1
	if next > 99 {
		return "", ErrMaxNotesPerDay
	}
	return fmt.Sprintf("%02d", next), nil
```

- [ ] **Step 2: Fix `findNoteByID` to use exact ID match (C7 complement)**

In `internal/core/utils.go`, replace line 73:
```go
			if !file.IsDir() && strings.HasPrefix(file.Name(), id+"_") {
```
With:
```go
			if !file.IsDir() {
				parts := strings.SplitN(file.Name(), "_", 2)
				if len(parts) > 0 && parts[0] == id {
```
And add a closing brace after the return statement on the next line. The full loop body becomes:
```go
		for _, file := range files {
			if !file.IsDir() {
				parts := strings.SplitN(file.Name(), "_", 2)
				if len(parts) > 0 && parts[0] == id {
					return filepath.Join(folderPath, file.Name()), nil
				}
			}
		}
```

- [ ] **Step 3: Fix `findNoteByID` to handle non-NotExist Stat errors (I4)**

In `internal/core/utils.go`, replace lines 47-49:
```go
	if _, err := os.Stat(BaseDir); os.IsNotExist(err) {
		return "", fmt.Errorf(constants.ErrNotesDirNotFound)
	}
```
With:
```go
	if _, err := os.Stat(BaseDir); err != nil {
		if os.IsNotExist(err) {
			return "", ErrNotesDirMissing
		}
		return "", fmt.Errorf("cannot access notes directory %s: %w", BaseDir, err)
	}
```

- [ ] **Step 4: Fix `findNoteByID` to surface permission errors on folders (I2)**

In `internal/core/utils.go`, replace lines 68-70:
```go
		files, err := os.ReadDir(folderPath)
		if err != nil {
			continue // Skip unreadable folders
		}
```
With:
```go
		files, err := os.ReadDir(folderPath)
		if err != nil {
			if os.IsPermission(err) {
				return "", fmt.Errorf("permission denied reading folder %s: %w", folderPath, err)
			}
			continue
		}
```

- [ ] **Step 5: Fix `List()` filepath.Rel error (I1)**

In `internal/core/list.go`, replace line 33:
```go
		rel, _ := filepath.Rel(BaseDir, path)
```
With:
```go
		rel, err := filepath.Rel(BaseDir, path)
		if err != nil {
			return fmt.Errorf("failed to compute relative path for %s: %w", path, err)
		}
```

- [ ] **Step 6: Fix `List()` Stat to handle non-NotExist errors (I4)**

In `internal/core/list.go`, replace lines 19-21:
```go
	if _, err := os.Stat(BaseDir); os.IsNotExist(err) {
		return notes, nil
	}
```
With:
```go
	if _, err := os.Stat(BaseDir); err != nil {
		if os.IsNotExist(err) {
			return notes, nil
		}
		return nil, fmt.Errorf("cannot access notes directory %s: %w", BaseDir, err)
	}
```

- [ ] **Step 7: Clean up unused constants import if needed**

Check if `internal/core/utils.go` still imports `constants` — it should no longer need `constants.ErrNoteNotFound` or `constants.ErrNotesDirNotFound` after Task 1. Remove the import if unused.

- [ ] **Step 8: Run tests**

Run: `go test ./... -v`
Expected: All pass. The error message changes are compatible with existing `strings.Contains(err.Error(), "not found")` checks because the sentinel `ErrNoteNotFound` message contains "not found".

- [ ] **Step 9: Run vet**

Run: `go vet ./...`
Expected: Clean

- [ ] **Step 10: Commit**

```
fix: core error handling — exact ID match, 99-note cap, permission errors (C7, I1, I2, I4)
```

---

### Task 5: Core Error Handling — edit.go, read.go, CLI (I3, I5, I9)

**Covers:** I3, I5, I9
**Files:**
- Modify: `internal/core/edit.go`
- Modify: `internal/core/read.go`
- Modify: `internal/cli/list.go:37`

- [ ] **Step 1: Fix `Edit` — named return + deferred Close error capture (I3, I5)**

Replace entire `internal/core/edit.go`:

```go
package core

import (
	"fmt"
	"os"

	"github.com/nickpricks/ft/internal/constants"
)

// Edit locates a note by its ID and appends the provided text to the bottom of the file.
func Edit(id string, text string) (path string, err error) {
	path, err = findNoteByID(id)
	if err != nil {
		return "", fmt.Errorf("failed to find note for editing: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, constants.FilePerm)
	if err != nil {
		return "", fmt.Errorf("failed to open note for editing: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close note file: %w", cerr)
		}
	}()

	if _, err := f.WriteString("\n" + text + "\n"); err != nil {
		return "", fmt.Errorf("failed to write to note: %w", err)
	}

	return path, nil
}
```

- [ ] **Step 2: Fix `Read` — wrap raw error (I5)**

Replace entire `internal/core/read.go`:

```go
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
```

- [ ] **Step 3: Fix `tabwriter.Flush` error in CLI list (I9)**

In `internal/cli/list.go`, replace line 37:
```go
	w.Flush()
	return nil
```
With:
```go
	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	return nil
```

- [ ] **Step 4: Run tests**

Run: `go test ./... -v`
Expected: All pass. The wrapped errors still contain "not found" so existing string checks work.

- [ ] **Step 5: Run vet**

Run: `go vet ./...`
Expected: Clean

- [ ] **Step 6: Commit**

```
fix: wrap raw errors in edit/read, capture Close errors, check Flush (I3, I5, I9)
```

---

### Task 6: Config Test Coverage (T1)

**Covers:** T1
**Files:**
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Create config test file**

```go
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/nickpricks/ft/internal/core"
)

// helper to set up config state for a test and restore it afterwards.
func setupConfigTest(t *testing.T) (tmpDir string, cleanup func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "ft-config-test-")
	if err != nil {
		t.Fatal(err)
	}

	origConfigPath := configFilePath
	origHomeDir := homeDir
	origInitErr := initErr
	origBaseDir := core.BaseDir

	cleanup = func() {
		configFilePath = origConfigPath
		homeDir = origHomeDir
		initErr = origInitErr
		core.BaseDir = origBaseDir
		os.RemoveAll(tmpDir)
	}

	homeDir = tmpDir
	initErr = nil
	configFilePath = filepath.Join(tmpDir, ".fmd.json")

	return tmpDir, cleanup
}

func TestLoadOrInit_ValidConfig(t *testing.T) {
	tmpDir, cleanup := setupConfigTest(t)
	defer cleanup()

	notesDir := filepath.Join(tmpDir, "my-notes")
	cfg := Config{NotesDir: notesDir}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		t.Fatal(err)
	}

	if err := LoadOrInit(); err != nil {
		t.Fatalf("LoadOrInit failed: %v", err)
	}

	if core.BaseDir != notesDir {
		t.Errorf("expected BaseDir=%q, got %q", notesDir, core.BaseDir)
	}
}

func TestLoadOrInit_CorruptJSON(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	if err := os.WriteFile(configFilePath, []byte("{broken json"), 0600); err != nil {
		t.Fatal(err)
	}

	err := LoadOrInit()
	if err == nil {
		t.Fatal("expected error for corrupt JSON, got nil")
	}
	if !containsStr(err.Error(), "invalid JSON") {
		t.Errorf("expected error mentioning 'invalid JSON', got: %v", err)
	}
}

func TestLoadOrInit_EmptyNotesDir(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	data, _ := json.Marshal(Config{NotesDir: ""})
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		t.Fatal(err)
	}

	err := LoadOrInit()
	if err == nil {
		t.Fatal("expected error for empty notes_dir, got nil")
	}
	if !containsStr(err.Error(), "empty notes_dir") {
		t.Errorf("expected error mentioning 'empty notes_dir', got: %v", err)
	}
}

func TestLoadOrInit_InitError(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	initErr = errors.New("simulated home dir failure")

	err := LoadOrInit()
	if err == nil {
		t.Fatal("expected error when initErr is set, got nil")
	}
	if !containsStr(err.Error(), "cannot load config") {
		t.Errorf("expected 'cannot load config' error, got: %v", err)
	}
}

func TestLoadOrInit_ConfigRoundtrip(t *testing.T) {
	tmpDir, cleanup := setupConfigTest(t)
	defer cleanup()

	notesDir := filepath.Join(tmpDir, "roundtrip-notes")
	cfg := Config{NotesDir: notesDir}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		t.Fatal(err)
	}

	// First load
	if err := LoadOrInit(); err != nil {
		t.Fatalf("first LoadOrInit failed: %v", err)
	}
	if core.BaseDir != notesDir {
		t.Errorf("first load: expected %q, got %q", notesDir, core.BaseDir)
	}

	// Reset and reload
	core.BaseDir = "notes"
	if err := LoadOrInit(); err != nil {
		t.Fatalf("second LoadOrInit failed: %v", err)
	}
	if core.BaseDir != notesDir {
		t.Errorf("second load: expected %q, got %q", notesDir, core.BaseDir)
	}
}

func TestLoadOrInit_PermissionDenied(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	// Write a valid config but make it unreadable
	data, _ := json.Marshal(Config{NotesDir: "/tmp/notes"})
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(configFilePath, 0000); err != nil {
		t.Skip("cannot change file permissions on this OS")
	}
	defer os.Chmod(configFilePath, 0600) // restore so cleanup can delete

	err := LoadOrInit()
	if err == nil {
		t.Fatal("expected error for unreadable config, got nil")
	}
	if !containsStr(err.Error(), "failed to read config file") {
		t.Errorf("expected 'failed to read config file' error, got: %v", err)
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && searchStr(s, substr)
}

func searchStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
```

Note: We use a custom `containsStr` helper instead of `strings.Contains` to keep it simple, but `strings.Contains` would also work fine — add the import if you prefer.

- [ ] **Step 2: Run config tests**

Run: `go test ./internal/config/ -v`
Expected: All 6 tests pass.

- [ ] **Step 3: Run full test suite**

Run: `go test ./... -v`
Expected: All pass

- [ ] **Step 4: Commit**

```
test: add config package tests — valid, corrupt JSON, empty dir, permissions (T1)
```

---

### Task 7: Core Test Improvements (T2-T6)

**Covers:** T2, T3, T4, T5, T6
**Files:**
- Modify: `internal/core/utils_test.go`
- Modify: `internal/core/edit_test.go`
- Modify: `internal/core/list_test.go`
- Modify: `internal/core/read_test.go`

- [ ] **Step 1: Add `TestFindNoteByID_MissingBaseDir` (T2) to utils_test.go**

Append to `internal/core/utils_test.go`:

```go
func TestFindNoteByID_MissingBaseDir(t *testing.T) {
	originalBaseDir := BaseDir
	BaseDir = "/nonexistent/path/that/does/not/exist"
	defer func() { BaseDir = originalBaseDir }()

	_, err := findNoteByID("01")
	if err == nil {
		t.Fatal("expected error when BaseDir doesn't exist, got nil")
	}
	if !errors.Is(err, ErrNotesDirMissing) {
		t.Errorf("expected ErrNotesDirMissing, got: %v", err)
	}
}
```

Add `"errors"` to the import block.

- [ ] **Step 2: Add `TestGetNextID_NonexistentDir` (T3) to utils_test.go**

Append to `internal/core/utils_test.go`:

```go
func TestGetNextID_NonexistentDir(t *testing.T) {
	id, err := getNextID("/nonexistent/dir")
	if err != nil {
		t.Fatalf("expected '01' for nonexistent dir, got error: %v", err)
	}
	if id != "01" {
		t.Errorf("expected '01', got %q", id)
	}
}
```

- [ ] **Step 3: Add `TestEdit_ReadOnlyFile` (T4) to edit_test.go**

Append to `internal/core/edit_test.go`:

```go
func TestEdit_ReadOnlyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-edit-readonly-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	originalBaseDir := BaseDir
	BaseDir = tmpDir
	defer func() { BaseDir = originalBaseDir }()

	dateFolder := GetDateFolder()
	if err := os.MkdirAll(dateFolder, 0755); err != nil {
		t.Fatal(err)
	}

	notePath := filepath.Join(dateFolder, "01_readonly.md")
	if err := os.WriteFile(notePath, []byte("read only note\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(notePath, 0444); err != nil {
		t.Skip("cannot change file permissions on this OS")
	}
	defer os.Chmod(notePath, 0644)

	_, err = Edit("01", "should fail")
	if err == nil {
		t.Error("expected error editing read-only file, got nil")
	}
}
```

- [ ] **Step 4: Add `TestList_MissingBaseDir` (T5) to list_test.go**

Append to `internal/core/list_test.go`:

```go
func TestList_MissingBaseDir(t *testing.T) {
	originalBaseDir := BaseDir
	BaseDir = "/nonexistent/path/that/does/not/exist"
	defer func() { BaseDir = originalBaseDir }()

	notes, err := List()
	if err != nil {
		t.Fatalf("expected empty slice for missing BaseDir, got error: %v", err)
	}
	if len(notes) != 0 {
		t.Errorf("expected 0 notes, got %d", len(notes))
	}
}
```

- [ ] **Step 5: Fix test setup error checking (T6) across all test files**

In `internal/core/edit_test.go`, change line 24:
```go
	os.MkdirAll(dateFolder, 0755)
```
To:
```go
	if err := os.MkdirAll(dateFolder, 0755); err != nil {
		t.Fatalf("test setup: %v", err)
	}
```

In `internal/core/read_test.go`, change line 24:
```go
	os.MkdirAll(dateFolder, 0755)
```
To:
```go
	if err := os.MkdirAll(dateFolder, 0755); err != nil {
		t.Fatalf("test setup: %v", err)
	}
```

In `internal/core/list_test.go`, change lines 35-43 — wrap each `os.MkdirAll` and `os.WriteFile` call:
```go
	if err := os.MkdirAll(filepath.Join(tmpDir, date1), 0755); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, date2), 0755); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, date1, "01_note_a.md"), []byte("A"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, date1, "02_note_b.md"), []byte("B"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, date2, "01_note_c.md"), []byte("C"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, date2, "02_ignore.txt"), []byte("Ignore"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
```

In `internal/core/utils_test.go`, change lines 46-48 — wrap each `os.WriteFile`:
```go
	if err := os.WriteFile(tmpDir+"/01_test.md", []byte("test"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.WriteFile(tmpDir+"/02_test.md", []byte("test"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.WriteFile(tmpDir+"/not_a_note.txt", []byte("test"), 0644); err != nil {
		t.Fatalf("test setup: %v", err)
	}
```

- [ ] **Step 6: Run all tests**

Run: `go test ./... -v`
Expected: All pass including the 4 new tests

- [ ] **Step 7: Run vet**

Run: `go vet ./...`
Expected: Clean

- [ ] **Step 8: Commit**

```
test: add edge case tests and fix setup error checking (T2-T6)
```

---

### Task 8: Documentation Fixes (D1, D4, D5)

**Covers:** D1, D4, D5 (D2 and D3 already done)
**Files:**
- Modify: `internal/core/types.go:4-6`
- Modify: `README.md:27-28`
- Modify: `cmd/feathertrailmd/README.md:11`

- [ ] **Step 1: Fix BaseDir comment (D1)**

In `internal/core/types.go`, replace lines 4-6:
```go
// BaseDir defines the root directory where notes are stored.
// In production, this defaults to "notes". During tests, it can be overridden.
var BaseDir = "notes"
```
With:
```go
// BaseDir defines the root directory where notes are stored.
// It is initialized to "notes" but is overridden at startup by
// config.LoadOrInit(), which reads from ~/.fmd.json or prompts
// the user. Tests also override it to use temporary directories.
var BaseDir = "notes"
```

- [ ] **Step 2: Fix README first-run example path (D4)**

In `README.md`, replace line 28:
```
Where would you like to store your notes? [C:\Users\username\Documents\FeatherTrailNotes]:
```
With:
```
Where would you like to store your notes? [~/Documents/FeatherTrailNotes]:
```

- [ ] **Step 3: Fix cmd README binary naming (D5)**

In `cmd/feathertrailmd/README.md`, replace line 11:
```
When you first run the `feathertrailmd` (or `ft`) binary, you'll be prompted to choose a configurable directory to store all your markdown notes in (defaults to `Documents/FeatherTrailNotes`).
```
With:
```
When you first run the binary, you'll be prompted to choose a directory to store your markdown notes in (defaults to `~/Documents/FeatherTrailNotes`).

> **Note:** `go install` produces a binary named `feathertrailmd`. The `ft` binary name is only produced by `make build`. You can alias it: `alias ft=feathertrailmd`
```

- [ ] **Step 4: Commit**

```
docs: fix BaseDir comment, README paths, binary naming clarification (D1, D4, D5)
```

---

### Task 9: NoteStore Refactor (A1, A2, A5)

**Covers:** A1, A2, A5
**Files:**
- Create: `internal/core/store.go`
- Create: `internal/core/store_test.go`
- Modify: `internal/core/add.go`
- Modify: `internal/core/list.go`
- Modify: `internal/core/read.go`
- Modify: `internal/core/edit.go`
- Modify: `internal/core/utils.go`
- Modify: `internal/core/types.go`
- Modify: `internal/config/config.go`
- Modify: `internal/cli/root.go`
- Modify: `internal/cli/add.go`
- Modify: `internal/cli/list.go`
- Modify: `internal/cli/read.go`
- Modify: `internal/cli/edit.go`
- Modify: All `internal/core/*_test.go`
- Modify: `internal/config/config_test.go`

This is the largest task. It replaces the global `BaseDir` with an encapsulated `NoteStore` struct.

- [ ] **Step 1: Write the failing store test**

Create `internal/core/store_test.go`:

```go
package core

import (
	"os"
	"strings"
	"testing"
)

func TestNewNoteStore_Valid(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-store-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatalf("NewNoteStore failed: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestNewNoteStore_EmptyDir(t *testing.T) {
	_, err := NewNoteStore("")
	if err == nil {
		t.Fatal("expected error for empty dir")
	}
}

func TestNoteStore_Add(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-store-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	path, err := store.Add("Test note from store")
	if err != nil {
		t.Fatalf("store.Add failed: %v", err)
	}
	if !strings.HasSuffix(path, "01_test_note_from_store.md") {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestNoteStore_ReadAfterAdd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-store-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	content := "Read this back"
	_, err = store.Add(content)
	if err != nil {
		t.Fatal(err)
	}

	got, err := store.Read("01")
	if err != nil {
		t.Fatalf("store.Read failed: %v", err)
	}
	if strings.TrimSpace(got) != content {
		t.Errorf("expected %q, got %q", content, strings.TrimSpace(got))
	}
}

func TestNoteStore_List(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-store-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Empty store
	notes, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 0 {
		t.Errorf("expected 0 notes, got %d", len(notes))
	}

	// Add and list
	store.Add("First note")
	store.Add("Second note")

	notes, err = store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 2 {
		t.Errorf("expected 2 notes, got %d", len(notes))
	}
}

func TestNoteStore_Edit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ft-store-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewNoteStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	store.Add("Original content")
	_, err = store.Edit("01", "Appended content")
	if err != nil {
		t.Fatalf("store.Edit failed: %v", err)
	}

	got, err := store.Read("01")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "Appended content") {
		t.Errorf("expected appended content, got: %s", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/core/ -run TestNewNoteStore -v`
Expected: FAIL — `NewNoteStore` not defined

- [ ] **Step 3: Create `internal/core/store.go`**

```go
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
```

- [ ] **Step 4: Move core functions to NoteStore methods**

Add methods to `internal/core/store.go` that delegate to the existing package-level functions after temporarily setting `BaseDir`. This is the incremental migration step — we keep the old functions working while adding the new API.

```go
// Add creates a new note in today's date folder.
func (s *NoteStore) Add(text string) (string, error) {
	prev := BaseDir
	BaseDir = s.baseDir
	defer func() { BaseDir = prev }()
	return Add(text)
}

// List returns all notes chronologically sorted.
func (s *NoteStore) List() ([]NoteInfo, error) {
	prev := BaseDir
	BaseDir = s.baseDir
	defer func() { BaseDir = prev }()
	return List()
}

// Read returns the raw content of a note by ID.
func (s *NoteStore) Read(id string) (string, error) {
	prev := BaseDir
	BaseDir = s.baseDir
	defer func() { BaseDir = prev }()
	return Read(id)
}

// Edit appends text to an existing note by ID.
func (s *NoteStore) Edit(id string, text string) (string, error) {
	prev := BaseDir
	BaseDir = s.baseDir
	defer func() { BaseDir = prev }()
	return Edit(id, text)
}
```

- [ ] **Step 5: Run store tests**

Run: `go test ./internal/core/ -run TestNoteStore -v`
Expected: All pass

- [ ] **Step 6: Update config to return value instead of mutating global (A2)**

Modify `internal/config/config.go`. Add a new function that returns the resolved directory:

```go
// ResolveNotesDir reads the config or prompts the user, returning the notes directory path.
// Unlike LoadOrInit, it does not mutate core.BaseDir.
func ResolveNotesDir() (string, error) {
	if initErr != nil {
		return "", fmt.Errorf("cannot load config: %w", initErr)
	}

	if isTestRun() {
		return "notes", nil
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to read config file %s: %w", configFilePath, err)
		}
	} else {
		var cfg Config
		if jsonErr := json.Unmarshal(data, &cfg); jsonErr != nil {
			return "", fmt.Errorf("config file %s contains invalid JSON: %w\nTo fix: delete the file and re-run ft", configFilePath, jsonErr)
		}
		if cfg.NotesDir == "" {
			return "", fmt.Errorf("config file %s has empty notes_dir field", configFilePath)
		}
		return cfg.NotesDir, nil
	}

	defaultDir := filepath.Join(homeDir, "Documents", "FeatherTrailNotes")

	fmt.Printf("Welcome to FeatherTrailMD!\n")
	fmt.Printf("It looks like this is your first time running the tool.\n")
	fmt.Printf("Where would you like to store your notes? [%s]: ", defaultDir)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input (is stdin available?): %w", err)
	}
	input = strings.TrimSpace(input)

	chosenDir := defaultDir
	if input != "" {
		chosenDir = input
	}

	chosenDir, err = filepath.Abs(chosenDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	cfg := Config{NotesDir: chosenDir}
	data, err = json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to save config: %w", err)
	}

	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		return "", fmt.Errorf("failed to write config file: %w", err)
	}

	if err := os.MkdirAll(chosenDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create notes directory: %w", err)
	}

	fmt.Printf("Awesome! Your notes will be saved in: %s\n\n", chosenDir)
	return chosenDir, nil
}

// LoadOrInit is kept for backwards compatibility. It calls ResolveNotesDir
// and sets core.BaseDir.
func LoadOrInit() error {
	dir, err := ResolveNotesDir()
	if err != nil {
		return err
	}
	core.BaseDir = dir
	return nil
}
```

- [ ] **Step 7: Wire NoteStore into CLI layer**

Modify `internal/cli/root.go` to create and hold a `NoteStore`:

```go
package cli

import (
	"github.com/nickpricks/ft/internal/config"
	"github.com/nickpricks/ft/internal/constants"
	"github.com/nickpricks/ft/internal/core"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     constants.RootUse,
		Short:   constants.RootShort,
		Long:    constants.RootLong,
		Example: constants.RootExample,
		Version: constants.Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "help" || cmd.CalledAs() == "help" {
				return nil
			}
			dir, err := config.ResolveNotesDir()
			if err != nil {
				return err
			}
			s, err := core.NewNoteStore(dir)
			if err != nil {
				return err
			}
			store = s
			return nil
		},
	}

	store *core.NoteStore
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Root flags can be added here
}
```

- [ ] **Step 8: Update CLI commands to use store**

Update `internal/cli/add.go`:
```go
package cli

import (
	"fmt"
	"strings"

	"github.com/nickpricks/ft/internal/constants"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     constants.AddUse,
	Short:   constants.AddShort,
	Long:    constants.AddLong,
	Example: constants.AddExample,
	Args:    cobra.MinimumNArgs(1),
	RunE:    runAdd,
}

func runAdd(cmd *cobra.Command, args []string) error {
	text := strings.Join(args, " ")
	path, err := store.Add(text)
	if err != nil {
		return err
	}
	fmt.Printf(constants.LogNoteCreated, path)
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
```

Update `internal/cli/list.go`:
```go
package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nickpricks/ft/internal/constants"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     constants.ListUse,
	Short:   constants.ListShort,
	Long:    constants.ListLong,
	Example: constants.ListExample,
	RunE:    runList,
}

func runList(cmd *cobra.Command, args []string) error {
	items, err := store.List()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Println(constants.LogNoNotes)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tID\tSLUG\tPATH")
	for _, note := range items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", note.Date, note.ID, note.Slug, note.Path)
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
```

Update `internal/cli/read.go`:
```go
package cli

import (
	"fmt"

	"github.com/nickpricks/ft/internal/constants"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:     constants.ReadUse,
	Short:   constants.ReadShort,
	Long:    constants.ReadLong,
	Example: constants.ReadExample,
	Args:    cobra.ExactArgs(1),
	RunE:    runRead,
}

func runRead(cmd *cobra.Command, args []string) error {
	id := args[0]
	content, err := store.Read(id)
	if err != nil {
		return err
	}
	fmt.Println(content)
	return nil
}

func init() {
	rootCmd.AddCommand(readCmd)
}
```

Update `internal/cli/edit.go`:
```go
package cli

import (
	"fmt"
	"strings"

	"github.com/nickpricks/ft/internal/constants"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     constants.EditUse,
	Short:   constants.EditShort,
	Long:    constants.EditLong,
	Example: constants.EditExample,
	Args:    cobra.MinimumNArgs(2),
	RunE:    runEdit,
}

func runEdit(cmd *cobra.Command, args []string) error {
	id := args[0]
	text := strings.Join(args[1:], " ")
	path, err := store.Edit(id, text)
	if err != nil {
		return err
	}
	fmt.Printf(constants.LogNoteUpdated, path)
	return nil
}

func init() {
	rootCmd.AddCommand(editCmd)
}
```

- [ ] **Step 9: Update config_test.go to test ResolveNotesDir**

Add to `internal/config/config_test.go`:

```go
func TestResolveNotesDir_ValidConfig(t *testing.T) {
	tmpDir, cleanup := setupConfigTest(t)
	defer cleanup()

	notesDir := filepath.Join(tmpDir, "resolve-notes")
	cfg := Config{NotesDir: notesDir}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		t.Fatal(err)
	}

	dir, err := ResolveNotesDir()
	if err != nil {
		t.Fatalf("ResolveNotesDir failed: %v", err)
	}
	if dir != notesDir {
		t.Errorf("expected %q, got %q", notesDir, dir)
	}
}
```

- [ ] **Step 10: Run full test suite**

Run: `go test ./... -v`
Expected: All pass

- [ ] **Step 11: Run vet**

Run: `go vet ./...`
Expected: Clean

- [ ] **Step 12: Commit**

```
refactor: introduce NoteStore struct, config returns value (A1, A2)
```

- [ ] **Step 13: Remove test detection heuristic (A5) — optional, deferred**

The `isTestRun()` heuristic in `config.go` still exists. Fully removing it (A5) requires making `LoadOrInit`/`ResolveNotesDir` accept an `io.Reader` parameter for stdin and a config path parameter. This is a larger refactor that can be done in a follow-up. For now, the NoteStore struct eliminates the need for `BaseDir` mutation in tests — tests create `NoteStore` instances directly with temp dirs.

If proceeding now: add `Options` struct to config as shown in CODE_REVIEW.md A5 section. Otherwise, leave this as a follow-up item.

---

### Task 10: NoteID Value Type (A3)

**Covers:** A3
**Files:**
- Create: `internal/core/noteid.go`
- Create: `internal/core/noteid_test.go`
- Modify: `internal/cli/read.go`
- Modify: `internal/cli/edit.go`

- [ ] **Step 1: Write the failing test**

Create `internal/core/noteid_test.go`:

```go
package core

import "testing"

func TestParseNoteID_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected NoteID
	}{
		{"1", NoteID("01")},
		{"01", NoteID("01")},
		{"5", NoteID("05")},
		{"12", NoteID("12")},
		{"99", NoteID("99")},
	}

	for _, tt := range tests {
		id, err := ParseNoteID(tt.input)
		if err != nil {
			t.Errorf("ParseNoteID(%q) error: %v", tt.input, err)
			continue
		}
		if id != tt.expected {
			t.Errorf("ParseNoteID(%q) = %q, want %q", tt.input, id, tt.expected)
		}
	}
}

func TestParseNoteID_Invalid(t *testing.T) {
	invalids := []string{"", "abc", "1a", "-1", "100", "0"}

	for _, input := range invalids {
		_, err := ParseNoteID(input)
		if err == nil {
			t.Errorf("ParseNoteID(%q) expected error, got nil", input)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/core/ -run TestParseNoteID -v`
Expected: FAIL — `NoteID` and `ParseNoteID` not defined

- [ ] **Step 3: Implement NoteID**

Create `internal/core/noteid.go`:

```go
package core

import (
	"fmt"
	"strconv"
)

// NoteID is a validated, zero-padded two-digit note identifier.
type NoteID string

// String returns the string representation of the NoteID.
func (id NoteID) String() string {
	return string(id)
}

// ParseNoteID validates and normalizes a raw ID string.
// Valid IDs are numeric strings representing 1-99, zero-padded to 2 digits.
func ParseNoteID(raw string) (NoteID, error) {
	n, err := strconv.Atoi(raw)
	if err != nil {
		return "", fmt.Errorf("invalid note ID %q: must be numeric", raw)
	}
	if n < 1 || n > 99 {
		return "", fmt.Errorf("invalid note ID %q: must be between 1 and 99", raw)
	}
	return NoteID(fmt.Sprintf("%02d", n)), nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/core/ -run TestParseNoteID -v`
Expected: All pass

- [ ] **Step 5: Wire NoteID into CLI read and edit commands**

Update `internal/cli/read.go` — add validation at the CLI boundary:

In the `runRead` function, change:
```go
	id := args[0]
	content, err := store.Read(id)
```
To:
```go
	noteID, err := core.ParseNoteID(args[0])
	if err != nil {
		return err
	}
	content, err := store.Read(noteID.String())
```

Add `"github.com/nickpricks/ft/internal/core"` to imports.

Update `internal/cli/edit.go` — same pattern:

In the `runEdit` function, change:
```go
	id := args[0]
	text := strings.Join(args[1:], " ")
	path, err := store.Edit(id, text)
```
To:
```go
	noteID, err := core.ParseNoteID(args[0])
	if err != nil {
		return err
	}
	text := strings.Join(args[1:], " ")
	path, err := store.Edit(noteID.String(), text)
```

Add `"github.com/nickpricks/ft/internal/core"` to imports.

- [ ] **Step 6: Run full test suite**

Run: `go test ./... -v`
Expected: All pass

- [ ] **Step 7: Run vet**

Run: `go vet ./...`
Expected: Clean

- [ ] **Step 8: Commit**

```
feat: add NoteID value type with validation at CLI boundary (A3)
```

---

## Final Step: Update CODE_REVIEW.md

After all tasks are complete, update `docs/CODE_REVIEW.md` — check off every resolved item (`[x]`). Items D2 and D3 were resolved before this plan. Item A5 may be partially resolved (NoteStore eliminates test-side need, but heuristic still exists in config). All other items should be fully resolved.

---

## Execution Summary

| Task | Items | Files Changed | Estimated Complexity |
|------|-------|---------------|---------------------|
| 1. Sentinel errors | A4 | 4 | Low |
| 2. Config error handling | C1-C6, I6, I7 | 1 | Medium |
| 3. Help/version skip | I8 | 1 | Low |
| 4. Core errors (utils/list) | C7, I1, I2, I4 | 2 | Medium |
| 5. Core errors (edit/read/cli) | I3, I5, I9 | 3 | Low |
| 6. Config tests | T1 | 1 (new) | Medium |
| 7. Core test improvements | T2-T6 | 4 | Low |
| 8. Documentation | D1, D4, D5 | 3 | Low |
| 9. NoteStore refactor | A1, A2, (A5) | ~15 | High |
| 10. NoteID type | A3 | 4 (2 new) | Low |
