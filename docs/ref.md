# FeatherTrailMD Codebase Reference (docs/ref.md)

This is a quick-glance reference guide to the FeatherTrailMD project structure, core functions, and constants.

## Directory Structure
- `cmd/ft/main.go`: The CLI entry point.
- `internal/cli/`: Cobra CLI commands (`root.go`, `add.go`, `list.go`, `read.go`, `edit.go`).
- `internal/notes/`: Pure filesystem logic (`add.go`, `list.go`, `read.go`, `edit.go`, `utils.go`, `types.go`).
- `internal/constants/constants.go`: Shared global constants for the entire app.
- `docs/`: Developer documentation and planning (`man.md`, `ref.md`, `PLAN.md`, `ActualPlan.md`).

## Core Functions (`internal/notes`)
- `Add(text string) (string, error)`: Creates a new note in today's folder.
- `List() ([]NoteInfo, error)`: Returns all notes chronologically sorted.
- `Read(id string) (string, error)`: Returns the string content of a note by ID.
- `Edit(id, text string) (string, error)`: Appends text to an existing note.
- `Slugify(text string) string`: Formats strings into clean filesystem slugs.
- `GetDateFolder() string`: Returns the target path for today's notes.
- `findNoteByID(id string) (string, error)`: Helps locate an absolute file path from a short ID.

## Core Constants (`internal/constants`)
- **Permissions**: `FilePerm (0644)`, `DirPerm (0755)`
- **Extensions**: `ExtMD (".md")`
- **Output Logs**: `LogNoteCreated`, `LogNoNotes`, `LogNoteUpdated`
- **Error Types**: `ErrNotesDirNotFound`, `ErrNoteNotFound` and formatting wrappers.
- **CLI Commands**: Standard Cobra Use/Short/Long strings for cohesive UI.
