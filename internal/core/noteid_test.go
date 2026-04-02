package core

import "testing"

func TestParseNoteID_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected NoteID
	}{
		{"1", NoteID("01")},
		{"01", NoteID("01")},
		{"5", NoteID("05")},
		{"12", NoteID("12")},
		{"99", NoteID("99")},
	}

	for _, tt := range tests {
		id, err := ParseNoteID(tt.input)
		if err != nil {
			t.Errorf("ParseNoteID(%q) error: %v", tt.input, err)
			continue
		}
		if id != tt.expected {
			t.Errorf("ParseNoteID(%q) = %q, want %q", tt.input, id, tt.expected)
		}
	}
}

func TestParseNoteID_Invalid(t *testing.T) {
	invalids := []string{"", "abc", "1a", "-1", "100", "0"}

	for _, input := range invalids {
		_, err := ParseNoteID(input)
		if err == nil {
			t.Errorf("ParseNoteID(%q) expected error, got nil", input)
		}
	}
}
