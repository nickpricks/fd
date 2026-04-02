# FeatherTrailMD Codebase Reference (docs/ref.md)

This is a quick-glance reference guide to the FeatherTrailMD project structure, core functions, and constants.

## Directory Structure
- `cmd/feathertrailmd/main.go`: The CLI entry point.
- `internal/cli/`: Cobra CLI commands (`root.go`, `add.go`, `list.go`, `read.go`, `edit.go`).
- `internal/core/`: Pure filesystem logic (`add.go`, `list.go`, `read.go`, `edit.go`, `utils.go`, `types.go`).
- `internal/config/config.go`: Global configuration package that manages `~/.fmd.json` and first-run setup.
- `internal/constants/constants.go`: Shared global constants for the entire app.
- `docs/`: Developer documentation and planning (`man.md`, `ref.md`, `PLAN.md`, `ActualPlan.md`).
- `Makefile`: Build automation (`make build`, `make test`, `make vet`, `make build-all`, etc.).

## Core Functions (`internal/core`)
- `Add(text string) (string, error)`: Creates a new note in today's folder.
- `List() ([]NoteInfo, error)`: Returns all notes chronologically sorted.
- `Read(id string) (string, error)`: Returns the string content of a note by ID.
- `Edit(id, text string) (string, error)`: Appends text to an existing note.
- `Slugify(text string) string`: Formats strings into clean filesystem slugs.
- `GetDateFolder() string`: Returns the target path for today's notes.
- `findNoteByID(id string) (string, error)`: Helps locate an absolute file path from a short ID.

## Config (`internal/config`)
- `Config` struct: Holds `NotesDir string` (the user's chosen notes storage directory).
- `LoadOrInit() error`: Reads `~/.fmd.json` or prompts the user on first run to choose a notes directory, then sets `core.BaseDir`.

## Core Constants (`internal/constants`)
- **Version**: `Version = "v0.1.6"` (used by the root Cobra command for `--version` output)
- **Permissions**: `FilePerm (0644)`, `DirPerm (0755)`
- **Extensions**: `ExtMD (".md")`
- **Output Logs**: `LogNoteCreated`, `LogNoNotes`, `LogNoteUpdated`
- **Error Types**: `ErrNotesDirNotFound`, `ErrNoteNotFound` and formatting wrappers.
- **CLI Commands**: All Cobra command strings are centralized here (`RootUse`, `RootShort`, `RootLong`, `AddUse`, `AddShort`, `ListUse`, `ListShort`, `ReadUse`, `ReadShort`, `EditUse`, `EditShort`, and their `Long`/`Example` variants).
