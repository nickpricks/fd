package core

import "errors"

var (
	ErrNoteNotFound    = errors.New("note not found")
	ErrNotesDirMissing = errors.New("notes directory not found")
	ErrCreateDateDir   = errors.New("failed to create date folder")
	ErrGenerateID      = errors.New("failed to generate next ID")
	ErrWriteNote       = errors.New("failed to write note")
	ErrMaxNotesPerDay  = errors.New("maximum notes per day (99) exceeded")
)
