package parse

import (
	"testing"
)

func TestEntryContent(t *testing.T) {
	var tests = []struct {
		explanation string
		input       string
		errExpected error
	}{
		{
			"standard entry content is valid",
			"hello, world!",
			nil,
		},
	}

	for _, tt := range tests {
		_, err := EntryContent(tt.input)
		if err != tt.errExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.input, err, tt.errExpected)
		}
	}
}
