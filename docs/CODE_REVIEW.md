# Code Review Report — FeatherTrailMD

**Date:** 2026-03-06
**Scope:** Full codebase (8 commits, 29 files, ~888 lines)
**Reviewed by:** 8 specialized agents (code quality, error handling, test coverage, type design, comments, silent failures)

---

## Quick Fix TODO List

Fixes are grouped by priority. Check them off as you go.

### Critical — Fix Before Next Release

- [x] **C1. Distinguish file-not-found from other ReadFile errors** (`internal/config/config.go:35-43`)
  Only treat `os.IsNotExist` as first-run. Return errors for permission denied, I/O failures, etc.

- [x] **C2. Surface JSON parse errors instead of silent re-prompt** (`internal/config/config.go:39-42`)
  If `~/.fmd.json` exists but contains invalid JSON, return an error telling the user their config is corrupt.

- [x] **C3. Stop silent fallback to relative `"notes"` path** (`internal/config/config.go:29-32`)
  When home directory is unknown, return an error instead of silently using CWD-relative `./notes`.

- [x] **C4. Check stdin read error** (`internal/config/config.go:60`)
  Replace `input, _ := reader.ReadString('\n')` with proper error checking. Fail in non-interactive environments.

- [x] **C5. Reuse home dir from init() instead of calling UserHomeDir again** (`internal/config/config.go:52-53`)
  Store `home` in a package variable during `init()` and reuse it, or check the error on the second call.

- [x] **C6. Fix Windows test detection** (`internal/config/config.go:46`)
  Check for `.test.exe` suffix and use `filepath.Separator` instead of hardcoded `/` in `/_go_build_`.

- [x] **C7. Fix ID prefix collision after 99 notes/day** (`internal/core/utils.go:110`)
  Either cap at 99 with an error, or change `findNoteByID` to exact-match the numeric prefix.

### Important — Should Fix Soon

- [x] **I1. Check `filepath.Rel` error in List()** (`internal/core/list.go:33`)
  Replace `rel, _ := filepath.Rel(...)` with error handling. Notes silently disappear from listings without this.

- [x] **I2. Handle unreadable folder errors in findNoteByID** (`internal/core/utils.go:69`)
  Don't silently `continue` on `ReadDir` errors. At minimum, surface permission errors.

- [x] **I3. Check `defer f.Close()` error in Edit** (`internal/core/edit.go:23`)
  Use a named return and deferred closure to capture close errors (data loss risk on NFS/disk-full).

- [x] **I4. Handle non-NotExist errors from `os.Stat`** (`internal/core/list.go:19`, `internal/core/utils.go:47`)
  Both `List()` and `findNoteByID` only check `os.IsNotExist`, ignoring permission errors.

- [x] **I5. Wrap raw OS errors with context** (`internal/core/edit.go:20-27`, `internal/core/read.go:17`)
  `edit.go` and `read.go` return bare OS errors. Wrap them like `add.go` does with `fmt.Errorf`.

- [x] **I6. Change config file permissions to 0600** (`internal/config/config.go:75`)
  `~/.fmd.json` is written world-readable (0644). Use 0600 for defense-in-depth.

- [x] **I7. Resolve user-supplied paths to absolute** (`internal/config/config.go:64`)
  Call `filepath.Abs(chosenDir)` before saving. Relative paths in the config cause CWD-dependent behavior.

- [x] **I8. Skip config init for --help and --version** (`internal/cli/root.go:15-17`)
  `PersistentPreRunE` runs on every command including help. First-run users see a prompt before seeing usage.

- [x] **I9. Check `w.Flush()` error in list CLI** (`internal/cli/list.go:37`)
  `tabwriter.Flush()` returns an error that is currently discarded.

### Tests — Add Coverage

- [x] **T1. Create `internal/config/config_test.go`** — 0% coverage on the package that controls where all notes are stored.
  - Test: valid config loads and sets `BaseDir`
  - Test: corrupt JSON returns descriptive error
  - Test: empty `NotesDir` returns error
  - Test: empty `configFilePath` returns error
  - Test: config file is written and can be re-read (roundtrip)

- [x] **T2. Test `findNoteByID` when BaseDir doesn't exist** (`internal/core/utils.go:47-49`)

- [x] **T3. Test `getNextID` on nonexistent directory returns "01"** (`internal/core/utils.go:85-90`)

- [x] **T4. Test `Edit` on a read-only file returns error** (`internal/core/edit.go`)

- [x] **T5. Test `List` when BaseDir doesn't exist returns empty slice** (`internal/core/list.go:19`)

- [x] **T6. Check errors on `os.MkdirAll`/`os.WriteFile` in test setup** (multiple test files)
  Use `t.Fatal(err)` guards instead of discarding errors.

### Documentation — Update Stale Content

- [x] **D1. Update `BaseDir` comment in types.go** (`internal/core/types.go:5`)
  Says "During tests, it can be overridden" — config now overrides it on every run.

- [x] **D2. Update code listing in docs/man.md** (`docs/man.md:40-64`)
  The `root.go` listing is missing the `PersistentPreRunE` hook.

- [x] **D3. Add `internal/config/` to docs/ref.md** (`docs/ref.md`)
  The reference guide doesn't mention the new config package.

- [x] **D4. Fix platform-specific path in README** (`README.md:27`)
  First-run example shows a Windows path. Show a platform-neutral placeholder or both OS examples.

- [x] **D5. Clarify binary naming in cmd README** (`cmd/feathertrailmd/README.md:11`)
  Says "(or ft)" but `go install` produces `feathertrailmd`, not `ft`. Only `make build` creates `ft`.

### Architecture — Longer Term

- [ ] **A1. Replace global `BaseDir` with a `NoteStore` struct**
  Eliminates global mutable state, enables parallel tests, makes dependencies explicit.

- [ ] **A2. Make config return a value instead of mutating a global**
  `LoadOrInit()` should return `(string, error)` instead of setting `core.BaseDir` as a side effect.

- [ ] **A3. Add a `NoteID` value type with validation**
  Reject non-numeric IDs at the CLI boundary instead of silently searching the filesystem.

- [ ] **A4. Add sentinel/typed errors**
  Replace string-matched errors with `var ErrNoteNotFound = errors.New(...)` for reliable error handling.

- [ ] **A5. Remove test-detection heuristic from production code**
  Replace `os.Args[0]` inspection with dependency injection (`io.Reader` for stdin, config path as parameter).

---

## Detailed Findings

### Section 1: Config Error Handling (Critical)

The `internal/config/config.go` file is the most concerning area. It uses a cascade of silent fallbacks where every error condition quietly degrades to default behavior. This creates a "silent data disassociation" risk: the tool appears to work, but notes are written to unexpected locations.

#### C1. ReadFile errors conflated with "file not found"

**File:** `internal/config/config.go`, lines 35-43

```go
data, err := os.ReadFile(configFilePath)
if err == nil {
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err == nil && cfg.NotesDir != "" {
        core.BaseDir = cfg.NotesDir
        return nil
    }
}
// ALL errors fall through to the first-run prompt
```

**Problem:** The `if err == nil` check means ANY `ReadFile` error — permission denied, I/O failure, disk corruption — is treated identically to "file does not exist." When this happens, the user is shown the first-run welcome wizard. If they press Enter, their existing config is overwritten on line 75.

**Scenario:** A user who has been using `ft` for months gets their `~/.fmd.json` permissions changed (e.g., by a backup tool). They run `ft list` and instead of seeing their notes, they get "Welcome to FeatherTrailMD!" Their old config is silently overwritten. All their notes in the original directory become invisible.

**Fix:**
```go
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
```

---

#### C2. JSON parse errors silently trigger re-prompt

**File:** `internal/config/config.go`, lines 39-42

```go
if err := json.Unmarshal(data, &cfg); err == nil && cfg.NotesDir != "" {
    core.BaseDir = cfg.NotesDir
    return nil
}
```

**Problem:** This compound conditional lumps three different failure modes into one silent fallthrough:
1. JSON parse error (corrupt file)
2. Valid JSON but wrong schema (e.g., `{"notes_directory": "/path"}` — note the typo)
3. Valid JSON but empty `NotesDir`

In all three cases, the user gets the first-run wizard with zero indication that their config was problematic. The corrupt file is then overwritten.

**Fix:** See the combined fix in C1 above — handle each case explicitly.

---

#### C3. Silent fallback to relative path

**File:** `internal/config/config.go`, lines 20-25 and 29-32

```go
func init() {
    home, err := os.UserHomeDir()
    if err == nil {
        configFilePath = filepath.Join(home, ".fmd.json")
    }
}

func LoadOrInit() error {
    if configFilePath == "" {
        core.BaseDir = "notes"  // <-- relative path, no warning
        return nil
    }
```

**Problem:** If `os.UserHomeDir()` fails in `init()`, `configFilePath` stays empty. Then `LoadOrInit()` silently sets `BaseDir` to the relative path `"notes"`. This means notes are written to `./notes` relative to whatever directory the user runs `ft` from. Running from different directories creates disconnected note stores. Notes appear "lost" when the user changes directories.

**Environments where this happens:** Docker containers without `$HOME`, cron jobs, systemd services, chrooted processes.

**Fix:**
```go
var homeDir string
var initErr error

func init() {
    home, err := os.UserHomeDir()
    if err != nil {
        initErr = fmt.Errorf("could not determine home directory: %w", err)
        return
    }
    homeDir = home
    configFilePath = filepath.Join(home, ".fmd.json")
}

func LoadOrInit() error {
    if initErr != nil {
        return fmt.Errorf("cannot load config: %w", initErr)
    }
    // ...
}
```

---

#### C4. Stdin read error discarded

**File:** `internal/config/config.go`, line 60

```go
input, _ := reader.ReadString('\n')
```

**Problem:** If stdin is closed (non-interactive env, piped input that ends), the error is discarded. `input` is empty, so `chosenDir` defaults to `defaultDir`. The config is silently saved with default settings — no human was ever asked.

**Fix:**
```go
input, err := reader.ReadString('\n')
if err != nil {
    return fmt.Errorf("failed to read input (is stdin available?): %w", err)
}
```

---

#### C5. Second UserHomeDir call discards error

**File:** `internal/config/config.go`, line 52

```go
home, _ := os.UserHomeDir()
defaultDir := filepath.Join(home, "Documents", "FeatherTrailNotes")
```

**Problem:** If this second call fails, `home` is `""` and `defaultDir` becomes the relative path `"Documents/FeatherTrailNotes"`. The user is shown this relative path as the default and may accept it. This code path is only reachable when the first `UserHomeDir()` call in `init()` succeeded but the second fails — unlikely but possible if `$HOME` is unset between calls.

**Fix:** Reuse the `homeDir` variable stored during `init()` (see C3 fix).

---

#### C6. Windows test detection failure

**File:** `internal/config/config.go`, line 46

```go
if strings.HasSuffix(os.Args[0], ".test") || strings.Contains(os.Args[0], "/_go_build_") {
```

**Problem:** On Windows:
- `go test` produces `.test.exe` binaries, not `.test`. The `HasSuffix` check fails.
- IDE-built test binaries use `\_go_build_` (backslash), not `/_go_build_`. The `Contains` check fails.

Result: on Windows, running `go test` on any package that triggers `LoadOrInit()` will hang waiting for stdin input.

**Fix:**
```go
exe := os.Args[0]
base := filepath.Base(exe)
if strings.HasSuffix(base, ".test") || strings.HasSuffix(base, ".test.exe") ||
    strings.Contains(exe, string(filepath.Separator)+"_go_build_") {
    core.BaseDir = "notes"
    return nil
}
```

Better long-term: remove the heuristic entirely (see A5) and use dependency injection.

---

### Section 2: Core Package Error Handling (Important)

#### C7. ID prefix collision after 99 notes/day

**File:** `internal/core/utils.go`, line 110

```go
return fmt.Sprintf("%02d", maxID+1), nil
```

**Problem:** `%02d` zero-pads to a minimum of 2 digits but has no maximum. After 99 notes in one day, the ID becomes `"100"`. The `findNoteByID` function uses prefix matching:

```go
if strings.HasPrefix(file.Name(), id+"_") {
```

Searching for ID `"10"` matches both `10_slug.md` and `100_slug.md`, returning whichever comes first.

**Fix (option A — cap at 99):**
```go
if maxID+1 > 99 {
    return "", fmt.Errorf("maximum notes per day (99) exceeded")
}
return fmt.Sprintf("%02d", maxID+1), nil
```

**Fix (option B — exact match):**
```go
// In findNoteByID, parse the numeric prefix and compare as integers
parts := strings.SplitN(file.Name(), "_", 2)
if parts[0] == id {
    // exact match
}
```

---

#### I1. filepath.Rel error discarded in List()

**File:** `internal/core/list.go`, line 33

```go
rel, _ := filepath.Rel(BaseDir, path)
```

**Problem:** If `filepath.Rel` fails, `rel` is empty. `strings.Split("", string(os.PathSeparator))` returns `[""]`, which has `len(parts) == 1 < 2`, so the note is silently dropped from the listing. The user sees an incomplete list.

**Fix:**
```go
rel, err := filepath.Rel(BaseDir, path)
if err != nil {
    return fmt.Errorf("failed to compute relative path for %s: %w", path, err)
}
```

---

#### I2. Unreadable folder errors suppressed in findNoteByID

**File:** `internal/core/utils.go`, line 69

```go
files, err := os.ReadDir(folderPath)
if err != nil {
    continue // Skip unreadable folders
}
```

**Problem:** If the note the user is looking for is in an unreadable folder, they get "note with ID X not found" — an actively misleading message that implies the note doesn't exist, when it actually exists but can't be accessed.

**Fix:**
```go
files, err := os.ReadDir(folderPath)
if err != nil {
    if os.IsPermission(err) {
        return "", fmt.Errorf("permission denied reading folder %s: %w", folderPath, err)
    }
    fmt.Fprintf(os.Stderr, "warning: skipping unreadable folder %s: %v\n", folderPath, err)
    continue
}
```

---

#### I3. defer f.Close() discards write errors in Edit

**File:** `internal/core/edit.go`, line 23

```go
defer f.Close()
```

**Problem:** On NFS or when disk is full, `Close()` can return errors because the actual write is deferred until close. The function returns success but data may not be persisted.

**Fix (requires named return):**
```go
func Edit(id string, text string) (path string, err error) {
    // ...
    f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, constants.FilePerm)
    if err != nil {
        return "", fmt.Errorf("failed to open note for editing: %w", err)
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = fmt.Errorf("failed to close note file: %w", cerr)
        }
    }()
    // ...
}
```

---

#### I4. os.Stat only checks IsNotExist

**Files:** `internal/core/list.go:19` and `internal/core/utils.go:47`

Both locations use:
```go
if _, err := os.Stat(BaseDir); os.IsNotExist(err) {
```

**Problem:** Permission denied or other `Stat` errors fall through silently. The subsequent `ReadDir` or `WalkDir` will fail with a confusing lower-level error.

**Fix (both locations):**
```go
if _, err := os.Stat(BaseDir); err != nil {
    if os.IsNotExist(err) {
        return notes, nil // or the appropriate "not found" error
    }
    return nil, fmt.Errorf("cannot access notes directory %s: %w", BaseDir, err)
}
```

---

#### I5. Raw OS errors returned without context

**Files:** `internal/core/edit.go:20-27` and `internal/core/read.go:17-19`

**Problem:** `edit.go` and `read.go` return bare OS errors (`return "", err`). Compare with `add.go` which wraps every error with `fmt.Errorf(constants.ErrWriteNote, err)`. The user sees `open /path/to/file: permission denied` instead of `failed to edit note: permission denied`.

**Fix:** Wrap consistently:
```go
// edit.go
return "", fmt.Errorf("failed to open note for editing: %w", err)

// read.go
return "", fmt.Errorf("failed to read note content: %w", err)
```

---

#### I6. Config file permissions

**File:** `internal/config/config.go`, line 75

```go
os.WriteFile(configFilePath, data, 0644)
```

**Fix:** Change to `0600`. The config currently only contains a directory path, but config files in home directories should be user-only readable by default. If tokens or credentials are ever added, 0644 would be a security vulnerability.

---

#### I7. User-supplied relative paths stored as-is

**File:** `internal/config/config.go`, lines 59-66

**Problem:** If a user types `./mynotes` at the prompt, it's stored literally. The config then behaves differently depending on CWD.

**Fix:**
```go
chosenDir, err := filepath.Abs(chosenDir)
if err != nil {
    return fmt.Errorf("failed to resolve path: %w", err)
}
```

---

#### I8. PersistentPreRunE runs on help/version

**File:** `internal/cli/root.go`, lines 15-17

**Problem:** A first-time user running `ft --help` gets an interactive setup prompt before seeing help text.

**Fix:**
```go
PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    if cmd.Name() == "help" || cmd.CalledAs() == "help" {
        return nil
    }
    return config.LoadOrInit()
},
```

Note: Cobra handles `--version` separately via `Version` field and does not trigger `PersistentPreRunE` for the `--version` flag specifically. But `ft version` as a subcommand would. Test to confirm behavior.

---

#### I9. tabwriter Flush error discarded

**File:** `internal/cli/list.go`, line 37

```go
w.Flush()
```

**Fix:**
```go
if err := w.Flush(); err != nil {
    return fmt.Errorf("failed to write output: %w", err)
}
```

---

### Section 3: Test Coverage

#### Current State

| Package | Coverage | Status |
|---------|----------|--------|
| `internal/core` | 84.8% | Good |
| `internal/config` | 0.0% | No tests exist |
| `internal/cli` | 0.0% | No tests exist |
| `cmd/feathertrailmd` | 0.0% | Acceptable (thin main) |
| `internal/constants` | N/A | No logic, no tests needed |

#### T1. Config package needs tests (Priority 1)

The config package is called on every single CLI invocation and controls where all data is stored. A regression here causes silent data loss. Create `internal/config/config_test.go` as a same-package test (`package config`) so it can manipulate `configFilePath` directly.

**Minimum viable tests:**

```go
func TestLoadOrInit_ValidConfig(t *testing.T) {
    // Write valid JSON to temp file, set configFilePath, verify core.BaseDir is set
}

func TestLoadOrInit_CorruptJSON(t *testing.T) {
    // Write "{broken" to config file, verify error returned (not silent fallthrough)
}

func TestLoadOrInit_EmptyNotesDir(t *testing.T) {
    // Write {"notes_dir": ""}, verify error returned
}

func TestLoadOrInit_EmptyConfigFilePath(t *testing.T) {
    // Set configFilePath = "", verify error returned
}

func TestLoadOrInit_ConfigRoundtrip(t *testing.T) {
    // Complete first-run flow with injected stdin, verify file is written and re-readable
}
```

**Testability blocker:** `LoadOrInit()` reads from `os.Stdin` directly. To test the interactive prompt path, refactor to accept `io.Reader`/`io.Writer`, or at minimum test only the non-interactive paths for now.

#### T2-T5. Core package coverage gaps

| Test | What it covers | File:Line |
|------|---------------|-----------|
| T2. `findNoteByID` with missing BaseDir | Error path when notes dir is deleted | `utils.go:47-49` |
| T3. `getNextID` on nonexistent dir | First-note-of-the-day creation | `utils.go:85-90` |
| T4. `Edit` on read-only file | Permission error propagation | `edit.go:20` |
| T5. `List` with missing BaseDir | Empty result on fresh install | `list.go:19` |

#### T6. Test setup error checking

Multiple test files discard errors from `os.MkdirAll` and `os.WriteFile` during setup. If setup fails, subsequent test assertions give confusing errors.

**Files affected:**
- `internal/core/edit_test.go:24`
- `internal/core/read_test.go:24`
- `internal/core/list_test.go:35-40`
- `internal/core/utils_test.go:46-48`

**Fix pattern:**
```go
// Before:
os.MkdirAll(dateFolder, 0755)

// After:
if err := os.MkdirAll(dateFolder, 0755); err != nil {
    t.Fatalf("test setup: %v", err)
}
```

#### Positive observations about existing tests

- Table-driven tests for `Slugify` — good Go practice
- `TestAdd_Fail` tests a real filesystem failure (BaseDir pointing to a file) — excellent
- `TestAdd_MultipleNotes` verifies ID incrementing — catches overwrite bugs
- `TestList` verifies sort order AND non-`.md` file filtering
- All tests use proper cleanup with `defer os.RemoveAll`

---

### Section 4: Documentation Issues

#### D1. Stale comment on BaseDir

**File:** `internal/core/types.go`, line 5

**Current:**
```go
// BaseDir defines the root directory where notes are stored.
// In production, this defaults to "notes". During tests, it can be overridden.
var BaseDir = "notes"
```

**Problem:** Since the config commit, `BaseDir` is overridden on every CLI invocation by `config.LoadOrInit()`. The comment implies only tests override it.

**Suggested:**
```go
// BaseDir defines the root directory where notes are stored.
// It is initialized to "notes" but is overridden at startup by
// config.LoadOrInit(), which reads from ~/.fmd.json or prompts
// the user. Tests also override it to use temporary directories.
var BaseDir = "notes"
```

#### D2. Stale code listing in docs/man.md

**File:** `docs/man.md`, lines 40-64

The inline code listing of `root.go` does not include the `PersistentPreRunE` hook or the `config` import. A new contributor reading this would have an incomplete picture of the CLI startup flow.

#### D3. Missing config package in docs/ref.md

**File:** `docs/ref.md`

The "quick-glance reference guide" lists the directory structure and core functions but doesn't mention `internal/config/` or `LoadOrInit()`.

#### D4. Windows-only path in README example

**File:** `README.md`, line 27

```
Where would you like to store your notes? [C:\Users\username\Documents\FeatherTrailNotes]:
```

The `go install` instruction above is bash-centric. Consider `[~/Documents/FeatherTrailNotes]` or showing both platforms.

#### D5. Binary naming ambiguity

**File:** `cmd/feathertrailmd/README.md`, line 11

Says `feathertrailmd (or ft)` but `go install` produces `feathertrailmd`. Only `make build` (which uses `-o ft`) creates an `ft` binary. Users who installed via `go install` would not find an `ft` command.

#### Additional comment improvements in config.go

| Line | Current | Suggestion |
|------|---------|------------|
| 30 | `// Fallback to local if home dir isn't found` | `// If the home directory cannot be determined, fall back to a relative "notes" directory in the current working directory.` |
| 38 | `// Config exists!` | Remove (noise) or replace with `// Parse and apply existing config` |
| 45-46 | `// If we're running under 'go test', bypass the prompt automatically.` | Add: `Note: The forward slash check may not match on Windows.` |
| 51 | `// config doesn't exist or is invalid` | `// config file is missing, unreadable, unparsable, or has an empty notes_dir` |

---

### Section 5: Type Design & Architecture

This section covers longer-term improvements. The current architecture works for v0.1.x but has structural issues that will compound as features are added.

#### A1. Replace global `BaseDir` with `NoteStore` struct

**Current design:**
```
config.LoadOrInit() --mutates--> core.BaseDir <--reads-- core.Add/List/Read/Edit
```

**Problems:**
- Exported mutable global with no encapsulation (any package can change it)
- No synchronization (data race if tests use `t.Parallel()` or any concurrency is added)
- Inverted dependency: `config` imports `core` just to set one variable
- Every test must manually save/override/restore with deferred cleanup

**Type design ratings for `BaseDir`:**
| Dimension | Score |
|-----------|-------|
| Encapsulation | 1/10 |
| Invariant Expression | 1/10 |
| Invariant Usefulness | 5/10 |
| Invariant Enforcement | 1/10 |

**Proposed:**
```go
type NoteStore struct {
    baseDir string
}

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

func (s *NoteStore) Add(text string) (string, error)          { ... }
func (s *NoteStore) Edit(id, text string) (string, error)     { ... }
func (s *NoteStore) Read(id string) (string, error)           { ... }
func (s *NoteStore) List() ([]NoteInfo, error)                { ... }
```

**Benefits:**
- `baseDir` is private — validated at construction, immutable after
- Tests create isolated instances with temp dirs — no global mutation
- Dependencies are explicit: CLI creates a `NoteStore` and passes it
- Safe for parallel tests and future concurrency
- Natural place for validation (non-empty, absolute path, writable)

#### A2. Config returns value instead of mutating global

**Current:** `LoadOrInit()` returns `error` and communicates its result by mutating `core.BaseDir`.

**Proposed:** `LoadOrInit()` returns `(string, error)`:
```go
func LoadOrInit() (string, error) {
    // ... resolution logic ...
    return resolvedDir, nil
}
```

The CLI layer uses the returned value to construct a `NoteStore`:
```go
PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    dir, err := config.LoadOrInit()
    if err != nil {
        return err
    }
    store, err = core.NewNoteStore(dir)
    return err
},
```

#### A3. NoteID value type

**Current:** Note IDs are passed as raw `string` from `args[0]` through to filesystem operations with no validation.

**Problem:** `ft read abc` silently searches the filesystem and returns "not found" instead of rejecting the invalid ID at the boundary.

**Proposed:**
```go
type NoteID string

func ParseNoteID(raw string) (NoteID, error) {
    if _, err := strconv.Atoi(raw); err != nil {
        return "", fmt.Errorf("invalid note ID %q: must be numeric", raw)
    }
    return NoteID(fmt.Sprintf("%02s", raw)), nil
}
```

#### A4. Sentinel/typed errors

**Current:** Errors are created with `fmt.Errorf` and string templates. Tests check errors with `strings.Contains`.

**Proposed:**
```go
var (
    ErrNoteNotFound    = errors.New("note not found")
    ErrNotesDirMissing = errors.New("notes directory not found")
)

// Enables: errors.Is(err, core.ErrNoteNotFound)
```

#### A5. Remove test detection heuristic

**Current:** Production code at `config.go:46` inspects `os.Args[0]` to detect test execution. This is fragile and bakes test awareness into production logic.

**Proposed:** Accept `io.Reader` for stdin input and config path as a parameter:
```go
type Options struct {
    ConfigPath string
    Stdin      io.Reader
    Stdout     io.Writer
}

func LoadOrInit(opts Options) (string, error) {
    // Uses opts.Stdin instead of os.Stdin
    // Uses opts.ConfigPath instead of package-level var
}
```

Tests pass `strings.NewReader("my/test/dir\n")` as stdin. No heuristic needed.

#### NoteInfo type design ratings

| Dimension | Score | Notes |
|-----------|-------|-------|
| Encapsulation | 2/10 | All fields exported, no constructor |
| Invariant Expression | 2/10 | Every field is `string` — Date not `time.Time`, ID not `int` |
| Invariant Usefulness | 6/10 | Constraints are real and sensible |
| Invariant Enforcement | 1/10 | No validation anywhere |

The struct is acceptable as a read-only display/transfer object for now. Adding a constructor that validates fields would be a small improvement. The `Date`-as-string works because ISO dates sort lexicographically, but this is coincidental and would break if the format changed.

#### Config type design ratings

| Dimension | Score | Notes |
|-----------|-------|-------|
| Encapsulation | 2/10 | Single exported field, package operates through globals |
| Invariant Expression | 2/10 | `NotesDir` is bare `string` |
| Invariant Usefulness | 4/10 | Minimal — one-field wrapper around a string |
| Invariant Enforcement | 3/10 | `LoadOrInit` checks for empty, but struct itself enforces nothing |

---

### Section 6: What's Good

This section documents the codebase's strengths so they are preserved as improvements are made.

1. **Clean package structure** — `cmd/`, `internal/cli/`, `internal/core/`, `internal/config/`, `internal/constants/` is idiomatic Go
2. **go vet passes cleanly** with no warnings
3. **Tests pass with race detector** (`go test -race ./...`)
4. **Core package at 84.8% coverage** — strong for a v0.1.x CLI tool
5. **Table-driven tests** for `Slugify` — follows Go best practice
6. **`TestAdd_Fail` tests a real failure** — `BaseDir` pointing to a file instead of a directory
7. **Proper `defer os.RemoveAll`** in all tests — no filesystem pollution
8. **Error wrapping in `add.go`** with `fmt.Errorf` and `%w` — the right pattern (needs consistency)
9. **Clean CLI/core separation** — cobra commands are thin wrappers around core functions
10. **Consistent formatting and naming** throughout
11. **`List()` sort verification** in tests catches regression in the descending date+ID ordering
12. **`List()` tests non-`.md` file filtering** — catches filter regressions
