package validate

import (
	"testing"
)

func TestEntryDate(t *testing.T) {
	var tests = []struct {
		explanation   string
		date          string
		validExpected bool
	}{
		{
			"standard date in 2019 is valid",
			"2019-10-18",
			true,
		},
		{
			"non-Friday date is invalid",
			"2019-10-19",
			false,
		},
		{
			"future date is invalid",
			"2039-03-13",
			false,
		},
		{
			"malformed date is invalid",
			"2019-10-1",
			false,
		},
		{
			"empty string is invalid",
			"",
			false,
		},
	}

	for _, tt := range tests {
		validActual := EntryDate(tt.date)
		if validActual != tt.validExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.date, validActual, tt.validExpected)
		}
	}
}
