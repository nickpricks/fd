# FeatherTrailMD (ft)

A blazing-fast, filesystem-first notes assistant built in Go.

## Overview
FeatherTrailMD (`ft`) is a markdown-based notes engine that organizes your notes cleanly into daily folders (`YYYY-MM-DD`). It uses simple commands to create, list, read, and append to your notes without any bloated databases or reliance on closed ecosystems.

## Installation

The easiest way to install FeatherTrailMD is via the `go` toolchain:

```bash
go install github.com/nickpricks/ft/cmd/feathertrailmd@latest
```

*(Note: The entry point was renamed from `cmd/ft` to `cmd/feathertrailmd` to adhere to Go standard layout while preventing `.gitignore` conflicts with the compiled `ft` binary.)*

Alternatively, you can run the provided installer scripts (`install.ps1` for Windows, `install.sh` for macOS/Linux).

## Usage

### First Run
When you run `ft` for the very first time, it will ask you where you'd like to store your Markdown notes globally:
```
Welcome to FeatherTrailMD!
It looks like this is your first time running the tool.
Where would you like to store your notes? [~/Documents/FeatherTrailNotes]:
```
Once configured, `ft` will save this path to `~/.fmd.json`, allowing you to take notes from any directory on your computer!

### Commands
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

- **Phase 1: Minimal CLI Notes Engine (🔄 In Progress)**: Core filesystem logic (`internal/core`), `add`, `list`, `read`, and `edit` commands with incremental day-based ID generation.
- **Phase 2: Metadata & Frontmatter (🔜 Planned)**: Custom line-by-line frontmatter parser to inject metadata and enable filtering via `list --status`, `--tag`, and `--date`.
- **Phase 3: Markdown Parser**: Custom Markdown parsing engine (Tokenizer, AST, HTML rendering).
- **Phase 4: Advanced AST Editing**: Intelligent content modification based on the AST (e.g. `edit --after`).
- **Phase 5: `mdcore` Library**: Extracting the formal parser into a reusable Go package.
- **Phase 6: Static Site Generator**: Converting notes into a blog with Go templates.
- **Phase 7: Multi-Deploy Targets**: Deploying directly to GitHub Pages, Cloudflare, etc.
- **Phase 8: Export & Extensions**: Exports (PDF, ePub) and advanced features (AI summaries, Graph views).
