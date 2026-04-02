# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

FeatherTrailMD (`ft`) — a filesystem-first Markdown notes CLI built in Go 1.26 with Cobra. Notes are stored as `<notesDir>/YYYY-MM-DD/ID_slug.md` where IDs are zero-padded incrementing integers per day (01, 02, ..., 99).

Currently at **Phase 1** (minimal CLI engine, complete). Phase 2 (frontmatter metadata & filtering) is next — see `docs/ActualPlan.md`.

## Build & Test Commands

```bash
make build        # Build the ft binary (cleans first)
make test         # Run all tests with -v
make vet          # Run go vet
make fmt          # Format code
make all          # tidy + fmt + vet + test + build
make cover        # Tests with coverage report → coverage.html
make install      # Build + install to GOPATH/bin
make build-all    # Cross-compile for linux/darwin/windows → bin/
```

Single test: `go test -run TestSlugify ./internal/core/`

## Architecture

```
cmd/feathertrailmd/main.go  → Entry point (calls cli.Execute)
internal/cli/               → Cobra command definitions (thin wrappers using NoteStore)
internal/core/              → NoteStore struct + unexported filesystem logic
internal/config/            → Config loading (~/.fmd.json), returns notes dir path
internal/constants/         → All string literals, permissions, version, CLI help text
```

**Why `cmd/feathertrailmd/` not `cmd/ft/`?** The compiled binary is named `ft`, which `.gitignore` excludes. Having source at `cmd/ft/` caused GitHub Actions to skip pushing the source code. The rename to `feathertrailmd` resolves this.

### Key patterns

- **`NoteStore` struct** (`store.go`): Encapsulates all note operations. Created in `PersistentPreRunE` from the config-resolved directory. Methods (`Add`, `List`, `Read`, `Edit`) delegate to unexported internal functions that accept `baseDir` as a parameter. No global mutable state.
- **Config returns a value**: `config.LoadOrInit()` returns `(string, error)` — the resolved notes directory. The CLI creates a `NoteStore` from it. Config does not import `core`.
- **CLI ↔ Core separation**: `cli/` commands call `store.Method()`. Business logic lives entirely in `core/` as unexported functions. `Slugify()` is the only exported utility (pure function, no state).
- **NoteID validation**: `ParseNoteID()` validates and normalizes IDs (numeric, 1-99, zero-padded) at the CLI boundary in `read` and `edit` commands.
- **Sentinel errors** (`errors.go`): `ErrNoteNotFound`, `ErrNotesDirMissing`, `ErrMaxNotesPerDay`, etc. Use `errors.Is()` for matching.
- **Note filenames**: `{ID}_{slug}.md` — the ID is exact-matched by `findNoteByID` across date folders (searches newest date first).
- **Test isolation**: Tests create a `NoteStore` with `NewNoteStore(tmpDir)` — no global save/restore needed. Config tests manipulate unexported package vars and set `os.Args[0]` to bypass `isTestRun()`.

## Version & Releases

Version string lives in `internal/constants/constants.go` (`Version` const). Releases are triggered by pushing a `v*` tag — GitHub Actions builds cross-platform binaries via `make build-all` and creates a GitHub Release (`.github/workflows/release.yml`).

## Remote

GitHub repo: `nickpricks/ft`
