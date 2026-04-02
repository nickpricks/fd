package core

import (
	"fmt"
	"strconv"
)

// NoteID is a validated, zero-padded two-digit note identifier.
type NoteID string

// String returns the string representation of the NoteID.
func (id NoteID) String() string {
	return string(id)
}

// ParseNoteID validates and normalizes a raw ID string.
// Valid IDs are numeric strings representing 1-99, zero-padded to 2 digits.
func ParseNoteID(raw string) (NoteID, error) {
	n, err := strconv.Atoi(raw)
	if err != nil {
		return "", fmt.Errorf("invalid note ID %q: must be numeric", raw)
	}
	if n < 1 || n > 99 {
		return "", fmt.Errorf("invalid note ID %q: must be between 1 and 99", raw)
	}
	return NoteID(fmt.Sprintf("%02d", n)), nil
}
