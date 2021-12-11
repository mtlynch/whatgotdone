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
		{
			"empty entry content is invalid",
			"",
			ErrEmptyEntryContent,
		},
		{
			"single-space content is invalid",
			" ",
			ErrEmptyEntryContent,
		},
		{
			"single-newline content is invalid",
			"\n",
			ErrEmptyEntryContent,
		},
		{
			"single carriage return content is invalid",
			"\r",
			ErrEmptyEntryContent,
		},
		{
			"tab content is invalid",
			"\t",
			ErrEmptyEntryContent,
		},
		{
			"whitespace-only content is invalid",
			"\t   \t\r\n \t",
			ErrEmptyEntryContent,
		},
	}

	for _, tt := range tests {
		_, err := EntryContent(tt.input)
		if err != tt.errExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.input, err, tt.errExpected)
		}
	}
}
