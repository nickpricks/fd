# Code Review Report — FeatherTrailMD

**Date:** 2026-03-06
**Scope:** Full codebase (8 commits, 29 files, ~888 lines)
**Reviewed by:** 8 specialized agents (code quality, error handling, test coverage, type design, comments, silent failures)
**Resolved:** 2026-04-02 on branch `feat/code-review-fixes`

---

## Quick Fix TODO List

### Critical — Fix Before Next Release

- [x] **C1. Distinguish file-not-found from other ReadFile errors** (`internal/config/config.go`)
- [x] **C2. Surface JSON parse errors instead of silent re-prompt** (`internal/config/config.go`)
- [x] **C3. Stop silent fallback to relative `"notes"` path** (`internal/config/config.go`)
- [x] **C4. Check stdin read error** (`internal/config/config.go`)
- [x] **C5. Reuse home dir from init() instead of calling UserHomeDir again** (`internal/config/config.go`)
- [x] **C6. Fix Windows test detection** (`internal/config/config.go`)
- [x] **C7. Fix ID prefix collision after 99 notes/day** (`internal/core/utils.go`)

### Important — Should Fix Soon

- [x] **I1. Check `filepath.Rel` error in List()** (`internal/core/list.go`)
- [x] **I2. Handle unreadable folder errors in findNoteByID** (`internal/core/utils.go`)
- [x] **I3. Check `defer f.Close()` error in Edit** (`internal/core/edit.go`)
- [x] **I4. Handle non-NotExist errors from `os.Stat`** (`internal/core/list.go`, `internal/core/utils.go`)
- [x] **I5. Wrap raw OS errors with context** (`internal/core/edit.go`, `internal/core/read.go`)
- [x] **I6. Change config file permissions to 0600** (`internal/config/config.go`)
- [x] **I7. Resolve user-supplied paths to absolute** (`internal/config/config.go`)
- [x] **I8. Skip config init for --help and --version** (`internal/cli/root.go`)
- [x] **I9. Check `w.Flush()` error in list CLI** (`internal/cli/list.go`)

### Tests — Add Coverage

- [x] **T1. Create `internal/config/config_test.go`**
- [x] **T2. Test `findNoteByID` when BaseDir doesn't exist**
- [x] **T3. Test `getNextID` on nonexistent directory returns "01"**
- [x] **T4. Test `Edit` on a read-only file returns error**
- [x] **T5. Test `List` when BaseDir doesn't exist returns empty slice**
- [x] **T6. Check errors on `os.MkdirAll`/`os.WriteFile` in test setup**

### Documentation — Update Stale Content

- [x] **D1. Update `BaseDir` comment in types.go**
- [x] **D2. Update code listing in docs/man.md**
- [x] **D3. Add `internal/config/` to docs/ref.md**
- [x] **D4. Fix platform-specific path in README**
- [x] **D5. Clarify binary naming in cmd README**

### Architecture

- [x] **A1. Replace global `BaseDir` with a `NoteStore` struct**
- [x] **A2. Make config return a value instead of mutating a global**
- [x] **A3. Add a `NoteID` value type with validation**
- [x] **A4. Add sentinel/typed errors**
- [ ] **A5. Remove test-detection heuristic from production code** — deferred; `isTestRun()` remains in config.go but is less critical now that NoteStore eliminates global state for core tests.

---

## What's Good

1. **Clean package structure** — `cmd/`, `internal/cli/`, `internal/core/`, `internal/config/`, `internal/constants/` is idiomatic Go
2. **go vet passes cleanly** with no warnings
3. **Tests pass with race detector** (`go test -race ./...`)
4. **Table-driven tests** for `Slugify` — follows Go best practice
5. **`TestAdd_Fail` tests a real failure** — baseDir pointing to a file instead of a directory
6. **Proper `defer os.RemoveAll`** in all tests — no filesystem pollution
7. **Consistent error wrapping** with `fmt.Errorf` and `%w` throughout
8. **Clean CLI/core separation** — Cobra commands are thin wrappers around `NoteStore` methods
9. **NoteStore encapsulation** — private `baseDir`, validated at construction, immutable after
10. **Config returns value** — no inverted dependency, no global mutation
11. **Sentinel errors** — `errors.Is()` support for reliable error matching
12. **NoteID boundary validation** — invalid IDs rejected at CLI layer before filesystem work
