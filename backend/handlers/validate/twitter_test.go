package validate

import (
	"testing"
)

func TestTwitterHandle(t *testing.T) {
	var tests = []struct {
		explanation   string
		handle        string
		validExpected bool
	}{
		{
			"well-formed handle address is valid",
			"jack",
			true,
		},
		{
			"single-character handle is valid",
			"j",
			true,
		},
		{
			"underscore characters are allowed",
			"jack_and_jill",
			true,
		},
		{
			"empty string is invalid",
			"",
			false,
		},
		{
			"handle with leading @ is invalid",
			"@jack",
			false,
		},
		{
			"handle with more than 15 characters is invalid",
			"jackandjillwentup",
			false,
		},
		{
			"handle with illegal characters is invalid",
			"jack.and.jill",
			false,
		},
	}

	for _, tt := range tests {
		validActual := TwitterHandle(tt.handle)
		if validActual != tt.validExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.handle, validActual, tt.validExpected)
		}
	}
}
