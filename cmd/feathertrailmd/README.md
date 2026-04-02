# feathertrailmd

This is the executable entry point for the FeatherTrailMD (`ft`) CLI notes tool.

## Installation

```bash
go install github.com/nickpricks/ft/cmd/feathertrailmd@latest
```

When you first run the binary, you'll be prompted to choose a directory to store your markdown notes in (defaults to `~/Documents/FeatherTrailNotes`).

> **Note:** `go install` produces a binary named `feathertrailmd`. The `ft` binary name is only produced by `make build`. You can alias it: `alias ft=feathertrailmd`

## Quick Start

```bash
ft add "My first note"
ft list
ft read 01
ft edit 01 "Appending some extra context"
```

For more details on the project architecture and the underlying `mdcore` parser, please see the [main repository README](../../README.md).
