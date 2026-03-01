# FeatherTrailMD (ft)

A blazing-fast, filesystem-first notes assistant built in Go.

## Overview
FeatherTrailMD (`ft`) is a markdown-based notes engine that organizes your notes cleanly into daily folders (`YYYY-MM-DD`). It uses simple commands to create, list, read, and append to your notes without any bloated databases or reliance on closed ecosystems.

## Installation

### Windows
Run the PowerShell installer script (requires administrative privileges if installing directly to GOPATH, or just use `make install` if you have `make` via Go/MinGW/etc.):
```powershell
.\install.ps1
```

### macOS / Linux
Run the bash installer script:
```bash
./install.sh
```

## Usage
`ft` is very straightforward to use:

- **Create a note**: `ft add "My note text here"`
- **List notes**: `ft list`
- **Read a note**: `ft read 01`
- **Append to a note**: `ft edit 01 "Appending some extra information"`

## Extending & Development
See the `docs/` folder for more information:
- `docs/man.md` - Developer manual and codebase execution flow.
- `docs/ref.md` - Quick reference to the codebase and constants.
- `docs/ActualPlan.md` - Current project roadmap and progression state.

## Roadmap & Vision
FeatherTrailMD is designed to naturally evolve from a simple CLI note tool into a comprehensive Markdown ecosystem. 

- **Phase 1: Minimal CLI Notes Engine (✅ Completed)**: Core filesystem logic (`internal/core`), `add`, `list`, `read`, and `edit` commands with incremental day-based ID generation.
- **Phase 2: Metadata & Frontmatter (🔜 Next Up)**: Custom line-by-line frontmatter parser to inject metadata and enable filtering via `list --status`, `--tag`, and `--date`.
- **Phase 3: Markdown Parser**: Custom Markdown parsing engine (Tokenizer, AST, HTML rendering).
- **Phase 4: Advanced AST Editing**: Intelligent content modification based on the AST (e.g. `edit --after`).
- **Phase 5: `mdcore` Library**: Extracting the formal parser into a reusable Go package.
- **Phase 6: Static Site Generator**: Converting notes into a blog with Go templates.
- **Phase 7: Multi-Deploy Targets**: Deploying directly to GitHub Pages, Cloudflare, etc.
- **Phase 8: Export & Extensions**: Exports (PDF, ePub) and advanced features (AI summaries, Graph views).
