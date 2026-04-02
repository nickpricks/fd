# Changelog

All notable changes to FeatherTrailMD are documented in this file.

---

## [v0.1.7] — 2026-04-02

Major internal quality release. Resolves 24 of 25 items from the [code review](docs/CODE_REVIEW.md). No breaking changes to CLI commands or note file format.

### Architecture

- **NoteStore struct** replaces global `BaseDir` — all note operations are now encapsulated in `core.NoteStore` with a private, validated `baseDir` field. No more global mutable state. ([A1])
- **Config returns value** — `config.LoadOrInit()` returns `(string, error)` instead of mutating a global. Config no longer imports `core`. ([A2])
- **Sentinel errors** — `ErrNoteNotFound`, `ErrNotesDirMissing`, `ErrMaxNotesPerDay`, etc. replace string-based error templates. Use `errors.Is()` for matching. ([A4])
- **NoteID type** — `ParseNoteID()` validates and normalizes note IDs (numeric, 1-99, zero-padded) at the CLI boundary. `ft read abc` now returns a clear error instead of silently searching. ([A3])

### Bug Fixes

- **Config: silent fallbacks eliminated** — permission errors, corrupt JSON, missing home dir, and stdin failures now return explicit errors instead of silently falling back to defaults. ([C1-C6])
- **Config file permissions** — `~/.fmd.json` is now written with `0600` instead of `0644`. ([I6])
- **Config paths resolved to absolute** — user-supplied relative paths are resolved via `filepath.Abs()` before saving. ([I7])
- **Help command skips config** — `ft --help` no longer triggers the first-run setup prompt. ([I8])
- **ID collision after 99 notes/day** — `getNextID` caps at 99 with a clear error. `findNoteByID` uses exact ID matching instead of prefix matching. ([C7])
- **Permission errors surfaced** — unreadable folders in `findNoteByID` and `List()` now return permission errors instead of silently skipping or returning misleading "not found". ([I2, I4])
- **filepath.Rel error handled** — `List()` no longer silently drops notes when `filepath.Rel` fails. ([I1])
- **Edit Close error captured** — `Edit()` uses named returns with deferred `Close()` to capture write errors on NFS/disk-full. ([I3])
- **Error wrapping** — `edit.go` and `read.go` now wrap OS errors with context, matching `add.go`'s pattern. ([I5])
- **tabwriter.Flush checked** — `list` CLI checks the `Flush()` error instead of discarding it. ([I9])
- **Windows test detection** — `isTestRun()` now checks `.test.exe` suffix and uses `filepath.Separator`. ([C6])

### Tests

- **Config package tests** — 6 new tests covering valid config, corrupt JSON, empty dir, init error, roundtrip, and permission denied. Config coverage went from 0% to meaningful. ([T1])
- **Edge case tests** — `findNoteByID` with missing dir, `getNextID` on nonexistent dir, `Edit` on read-only file, `List` with missing dir. ([T2-T5])
- **Setup error checking** — all test files now check `os.MkdirAll`/`os.WriteFile` errors with `t.Fatal()` guards. ([T6])
- **NoteStore tests** — constructor validation (empty dir, relative dir) and integration test (Add → Read → List → Edit cycle).
- **Test count**: 9 → 25

### Documentation

- Updated `docs/man.md` — current code listings for `root.go`, `add.go`, new config section. ([D2])
- Updated `docs/ref.md` — added config package, Makefile, Version constant. ([D3])
- Updated `docs/PLAN.md` — Phase 1 marked ✅, `fmd` → `ft`, repo structure updated. 
- Updated `docs/ActualPlan.md` — release verified, config subsection added.
- Fixed `README.md` — cross-platform path in first-run example. ([D4])
- Fixed `cmd/feathertrailmd/README.md` — binary naming clarification. ([D5])
- Fixed `types.go` BaseDir comment. ([D1])
- Cleaned `CODE_REVIEW.md` — removed example code, kept checklist only.
- Added `CLAUDE.md` for Claude Code guidance.

---

## [v0.1.6] — 2026-03-06

- Added global config storage (`~/.fmd.json`) with first-run prompt.
- Added `internal/config/` package.
- Added `internal/constants/` for centralized strings.
- Cross-platform installer scripts.

## [v0.1.5] — 2026-03-06

- Renamed `cmd/ft` to `cmd/feathertrailmd` to resolve `.gitignore` conflict.

## [v0.1.4] — 2026-03-06

- Initial release with `add`, `list`, `read`, `edit` commands.
- Per-day incremental ID generation.
- Makefile with build, test, vet, cross-compile targets.
- GitHub Actions release workflow.
