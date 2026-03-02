package constants

import "os"

// File system constants
const (
	FilePerm os.FileMode = 0644
	DirPerm  os.FileMode = 0755
	ExtMD                = ".md"
)

// App Constants
const (
	Version = "v0.1.5"
)

// Log messages
const (
	LogNoteCreated = "Note created: %s\n"
	LogNoteUpdated = "Note updated: %s\n"
	LogNoNotes     = "No notes found."
)

// Error string templates
const (
	ErrNotesDirNotFound = "notes directory not found"
	ErrNoteNotFound     = "note with ID %s not found"
	ErrCreateDateFolder = "failed to create date folder: %w"
	ErrGenerateID       = "failed to generate next ID: %w"
	ErrWriteNote        = "failed to write note: %w"
)

// Root command texts
const (
	RootUse   = "ft"
	RootShort = "FeatherTrailMD is a quick notes tool - Super Fast Thoughts Notes"
	RootLong  = `FeatherTrailMD (ft) is a simple, filesystem-first notes assistant.

To install or upgrade this CLI tool:
- Windows: Run .\install.ps1 or .\update.ps1
- macOS/Linux: Run ./install.sh or ./update.sh`
	RootExample = `  ft add "Hello World!"
  ft list
  ft read 01`
)

// Add command texts
const (
	AddUse     = "add [text...]"
	AddShort   = "Quickly capture a new note"
	AddLong    = "The add command takes any text you provide and creates a new markdown note in today's folder. It automatically generates an incremental ID and a smart slug based on your content."
	AddExample = `  ft add "Meeting notes with the engineering team"
  ft add Validate the new API endpoints`
)

// List command texts
const (
	ListUse     = "list"
	ListShort   = "Display all notes chronologically"
	ListLong    = "The list command scans your notes directory and outputs a formatted table of all your notes, sorted from newest to oldest."
	ListExample = `  ft list`
)

// Read command texts
const (
	ReadUse     = "read [id]"
	ReadShort   = "Read the contents of a specific note"
	ReadLong    = "The read command outputs the full markdown content of a note to the terminal. It searches for the note ID in today's folder first, then searches backwards through recent days."
	ReadExample = `  ft read 01
  ft read 12`
)

// Edit command texts
const (
	EditUse     = "edit [id] [text...]"
	EditShort   = "Append text to the bottom of an existing note"
	EditLong    = "The edit command allows you to add more information to a note you've already created. In this version, editing purely appends the new text to the bottom of the file."
	EditExample = `  ft edit 01 "Adding another point to the meeting notes"`
)
