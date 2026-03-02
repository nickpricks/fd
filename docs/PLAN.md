# 🪶 FeatherTrailMD (fmd): The Dreams (PLAN.md)

> "A feather-light trail of thoughts, captured in Markdown, evolving into a digital forest."

**FeatherTrailMD** (`fmd`) is a terminal-first Markdown system that evolves in three layers:
1. A fast CLI note tool (filesystem based)
2. A reusable Markdown parsing library (`mdcore`)
3. A static site generator with deploy and export capabilities

## 💎 Primary Philosophy
* **Filesystem first**: The source of truth is your local files.
* **Markdown canonical**: Data is always readable and portable.
* **No database**: Keep it light, keep it simple.
* **Git friendly**: Versioning is built into the workflow.
* **Modular from day one**: Designed to grow without bloat.
* **Portable**: Single binary, minimal configuration.

---

# 📂 Repository Structure (Target State)
```
/cmd/fmd          # CLI Entry Point
/internal/cli     # CLI Command routing
/internal/core   # Note management logic
/internal/site    # SSG logic
/internal/deploy  # Deployment providers
/internal/export  # Export engines (PDF, etc)
/pkg/mdcore       # The core parser library
```

---

# 🚀 The 8-Phase Roadmap

## Phase 1 – Minimal CLI Notes Engine
**Goal**: Build a working note tool with per-day incremental IDs.
* **File Structure**:
  ```
  notes/
    YYYY-MM-DD/
      01_slug.md
      02_slug.md
  ```
* **Features**:
  * `ft add "text"`
  * `ft list`
  * `ft read 01`
  * `ft edit 01 "append text"`
  * `ft done 01` (Note: Status tracking moved to Phase 2)
* **Rules**: Per-day incremental numbering, Slug from first 3–5 words, Pure append, No frontmatter yet.
* **Infrastructure**: Test coverage, `Makefile`, and developer manual (`docs/man.md`).
* **Packages & releases**: Github actions maybae

## Phase 2 – Metadata + Frontmatter
**Goal**: Add structured metadata.
* **Format**: Custom Rules or JSON within `---` (e.g., `status: draft` or `{"status": "draft"}`).
  ```
  ---
  status: draft
  created: 2026-03-01T12:30:00
  tags: go,cli
  ---
  ```
* **New Features**:
  * `fmd done 01`
  * `fmd list --status draft`
  * `fmd list --tag go`
  * `fmd list --date 2026-03-01`
* **Implementation**: Manual extraction and parsing (No YAML engine).

## Phase 3 – Markdown Parser (Learning Phase)
**Goal**: Build a custom Markdown parsing engine.
* **Scope**: Tokenizer, Block parser (Headers, Lists, Paragraphs, Code blocks), AST representation, Markdown & HTML renderers.
* **CLI Additions**:
  * `fmd preview 01`
  * `fmd export 01 --html`

## Phase 4 – Advanced Editing + AST Manipulation
**Goal**: Intelligent content modification via AST.
* **Features**:
  * `fmd edit 01 "text" --after "XYZ"`
  * `fmd edit 01 "text" --before "ABC"`
  * `fmd edit 01 --replace "old" "new"`
  * Task checkbox parsing, Link extraction, Directive block support.

## Phase 5 – Extract `mdcore` Library
**Goal**: Move parser code from `/internal/parser` to `/pkg/mdcore` and formalize the Public API.

## Phase 6 – Static Site Generator
**Goal**: Convert notes into a professional blog using Go templates.
* **Features**: `fmd publish 01`, `fmd build`, `fmd serve`.
* **Output**: Index, Post, Tag pages, RSS, Sitemap.

## Phase 7 – Multi Deploy Targets
**Goal**: One-command deployment via abstract `Deployer` interface.
* **Targets**: GitHub Pages, Cloudflare, Netlify.

## Phase 8 – Export & Extensions (The Horizon)
**Goal**: Modular export layer and AI-powered features.
* **Export**: PDF, ePub, ZIP, Google Drive.
* **Extensions**: AI auto-summary, Graph view (Obsidian-style), Git history explorer, Custom shortcodes, Plugin marketplace.

---

# 🏁 MVP Definition
**MVP = Phase 1 + Phase 2**
* Deliverables: Add, List, Read, Edit (append), Done, Metadata filtering.

---

# ⏳ Estimated Timeline
| Phase | Nicks ETA | AntiGravity Estimate |
| :--- | :--- | :--- |
| 1 | 10 Hours | Solidly achievable |
| 2 | 1 week | High confidence |
| 3 | 2 weeks | Complex (custom parser) |
| 4 | 2 weeks | Requires robust AST |
| 5 | 1 week | Refactoring focus |
| 6 | 2 weeks | Template work |
| 7 | 2 weeks | API integrations |
| 8 | Open ended | Infinite potential 🚀 |

**Total Roadmap: ~11 weeks.**
