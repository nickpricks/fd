# WORKPLAN

## Current Focus: Code Review Resolution (Phase 1 Polish)
Based on `docs/CODE_REVIEW.md` (2026-03-06).

### ⏳ Time Estimates

| Category | Nick's ETA | AntiGravity's ETA |
| :--- | :--- | :--- |
| **🔴 Critical Fixes** (C1-C7) | ~1.5 hours | ~15 mins |
| **🟡 Important Fixes** (I1-I9) | ~1.5 hours | ~20 mins |
| **🔵 Tests** (T1-T6) | ~2.5 hours | ~25 mins |
| **🟢 Documentation** (D1-D5) | ~45 mins | ~10 mins |
| **Total Target Sweep** | **~6 hours** | **~1 hour 10 mins** |

---

### 🔴 Critical Fixes
- [x] **C1. Distinguish file-not-found from other ReadFile errors** (`internal/config/config.go`)
- [x] **C2. Surface JSON parse errors instead of silent re-prompt** (`internal/config/config.go`)
- [x] **C3. Stop silent fallback to relative `"notes"` path** (`internal/config/config.go`)
- [x] **C4. Check stdin read error** (`internal/config/config.go`)
- [x] **C5. Reuse home dir from init() instead of calling UserHomeDir again** (`internal/config/config.go`)
- [x] **C6. Fix Windows test detection** (`internal/config/config.go`)
- [x] **C7. Fix ID prefix collision after 99 notes/day** (`internal/core/utils.go`)

### 🟡 Important Fixes
- [x] **I1. Check `filepath.Rel` error in List()** (`internal/core/list.go`)
- [x] **I2. Handle unreadable folder errors in findNoteByID** (`internal/core/utils.go`)
- [x] **I3. Check `defer f.Close()` error in Edit** (`internal/core/edit.go`)
- [x] **I4. Handle non-NotExist errors from `os.Stat`** (`internal/core/list.go`, `internal/core/utils.go`)
- [x] **I5. Wrap raw OS errors with context** (`internal/core/edit.go`, `internal/core/read.go`)
- [x] **I6. Change config file permissions to 0600** (`internal/config/config.go`)
- [x] **I7. Resolve user-supplied paths to absolute** (`internal/config/config.go`)
- [x] **I8. Skip config init for --help and --version** (`internal/cli/root.go`)
- [x] **I9. Check `w.Flush()` error in list CLI** (`internal/cli/list.go`)

### 🔵 Tests
- [x] **T1. Create `internal/config/config_test.go`**
- [x] **T2. Test `findNoteByID` when BaseDir doesn't exist**
- [x] **T3. Test `getNextID` on nonexistent directory returns "01"**
- [x] **T4. Test `Edit` on a read-only file returns error**
- [x] **T5. Test `List` when BaseDir doesn't exist returns empty slice**
- [x] **T6. Check errors on `os.MkdirAll`/`os.WriteFile` in test setup**

### 🟢 Documentation
- [x] **D1. Update `BaseDir` comment in types.go**
- [x] **D2. Update code listing in docs/man.md**
- [x] **D3. Add `internal/config/` to docs/ref.md**
- [x] **D4. Fix platform-specific path in README**
- [x] **D5. Clarify binary naming in cmd README**

### ⚪ Architecture (Push to Phase 2 or Next Major Refactor)
- [ ] A1-A5: Deferred for future structural improvements.
