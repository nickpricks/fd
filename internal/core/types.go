// Package notes provides the core filesystem logic for FeatherTrailMD.
package core

// BaseDir defines the root directory where notes are stored.
// In production, this defaults to "notes". During tests, it can be overridden.
var BaseDir = "notes"

// NoteInfo represents the metadata of a single note extracted from the filesystem.
type NoteInfo struct {
	Path string // Full relative path to the file
	Date string // Date string extracted from the parent folder name
	ID   string // The incremental ID (e.g., "01")
	Slug string // The text slug of the note
}
