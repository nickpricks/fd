# 🛠️ FeatherTrailMD Actual Doable Flow (ActualPlan.md)

This document provides the technical, point-by-point "how-to" for the first two phases of FeatherTrailMD (`ft`).

## Phase 1: Minimal CLI Notes Engine (COMPLETED ✅)
**Goal**: Build a working note tool with per-day incremental IDs and slug generation.

### 🔌 1. CLI Entry Point (`cmd/ft/main.go`)
- **Technology**: Use `spf13/cobra` for the CLI framework.
- **Root**: `ft`
- **Commands**:
    - `ft add [text]`: Creates a new note.
    - `ft list`: Lists all notes in chronological order.
    - `ft read [ID]`: Displays the content of a specific note.
    - `ft edit [ID] [text]`: Appends text to an existing note.

### 📂 2. Filesystem Orchestration (`internal/core/`)
- **Modularity**: The monolithic `notes.go` was split into domain-specific files (`add.go`, `list.go`, `read.go`, `edit.go`, `utils.go`, `types.go`) for maintainability.
- **Base Directory**: Default to `./notes` in the current working directory.
- **Structure**: `notes/YYYY-MM-DD/ID_slug.md`.
- **ID Logic**: 
    1. Check if folder `notes/YYYY-MM-DD/` exists.
    2. Read files in folder using `os.ReadDir()`.
    3. Find the highest prefix (e.g., `01`, `02`).
    4. Next ID = Highest + 1 (padded to 2 digits).
- **Slug Logic**: 
    1. Take the first 3–5 words of the input text.
    2. Lowercase and replace spaces/special characters with underscores.

### 🧪 3. Testing Infrastructure
- **Location**: Found in `internal/core/*_test.go`
- **Coverage**: Add, Edit, Read, List commands, plus string utility functions.
- **Strategy**: Creating temporary, isolated `notes` directories (`os.MkdirTemp`) during tests to prevent polluting actual user data.

### ⚙️ 4. Build & Automation (`Makefile`)
- **Location**: Root directory.
- **Commands**: `make all`, `make clean`, `make test`, `make build`, etc., to streamline the development loop.
- **Distribution & Updates**: Added `make build-all` for cross-platform binary compilation (Linux, macOS, Windows), `make install` to install locally to GOPATH, and `make upgrade` to auto-update dependencies and reinstall.

### 📚 5. Documentation
- **Developer Manual**: Added `docs/man.md` featuring a line-by-line code explanation of how Cobra commands (`main`, `root`, `add_test`) are wired together for new contributors.
- [x] **Developer Manual Continued** - Add, Edit, Read, List commands.
- [x] **Codebase Reference**: Added `docs/ref.md` to serve as a quick reference for the project's structure, key functions, and constants.

### 🌐 6. Cross-Platform TODOs
- [x] Audit `filepath.Join` and `filepath.ToSlash` across OSes, especially in `internal/core/list.go`, to ensure flawless Linux/Mac compatibility.
- [x] Ensure terminal coloring (if added in the future) relies on cross-platform libraries (e.g. `fatih/color`).
- [x] Ensure help has install/update commands information
- [x] build and installer scripts, windows, linux & macos
- [x] use installer to figure out any issues, if any update @ActualPlan.md
- [x] run installer update command to fetch latest version if available

### GIT
- [x] Initialize git repository, prep first commit
- [x] Add & verify .gitignore
- [ ] Add github (or any other) remote
- [x] Add & verify initial README.md
- [x] Commit & push initial code

---

## Phase 2: Metadata & Frontmatter (NEXT UP 🔜)
**Update TESTS and verify PH 1 changes**
**Goal**: Add structured metadata and filtering support natively (No YAML dependencies).

### 📑 1. Frontmatter Format
Every note will eventually have a custom block at the top:
```
---
status: draft
created: 2026-03-01T14:22:00
tags: go,cli
---
```

### 🛠️ 2. Implementation Points
- **Manual Parser (`internal/core/meta.go`)**: 
    - Implement a custom line-by-line parser to extract key-value pairs (`key: value`) between `---` markers.
    - Avoid `yaml.v3` to keep it lightweight!
- **Auto-Injection**: 
    - When `ft add` is called, automatically prepend the initial frontmatter block before writing the user's text.
- **Command Updates**:
    - **`ft done [ID]`**: Reads the file, parses the custom metadata block, changes `status: draft` to `status: done`, and rewrites the file.
- **Filtering Logic**: 
    - Update `ft list` to accept flags: `--status`, `--tag`, `--date`.
    - Logic: For each file found during listing, read the first few lines, extract the custom metadata block, and apply the flag filters.

---

## 🏗️ Dev Setup & Project State
- Module initialized as: `github.com/nickpricks/ft`
- Dependencies: `github.com/spf13/cobra` (and no others yet!)
- Phase 1 Code is living in `cmd/ft/` and `internal/cli/`, `internal/core/`.
