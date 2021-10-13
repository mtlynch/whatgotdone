package parse

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
			"future Friday date is invalid",
			"2039-03-11",
			false,
		},
		{
			"date prior to What Got Done epoch is invalid",
			"2018-12-28",
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
		_, err := EntryDate(tt.date)
		if (err == nil) != tt.validExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.date, err, tt.validExpected)
		}
	}
}
