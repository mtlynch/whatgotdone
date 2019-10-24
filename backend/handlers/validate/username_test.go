package validate

import (
	"strings"
	"testing"
)

func TestUsername(t *testing.T) {
	var tests = []struct {
		explanation   string
		username      string
		validExpected bool
	}{
		{
			"well-formed handle address is valid",
			"mike",
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
			"handle with exactly 60 characters is valid",
			strings.Repeat("A", 60),
			true,
		},
		{
			"empty string is invalid",
			"",
			false,
		},
		{
			"'undefined' as a username is invalid",
			"undefined",
			false,
		},
		{
			"Asterisk as a username is invalid because it has different interpretation in DBs like Redis",
			"*",
			false,
		},
		{
			"Asterisk in the middle of a username is invalid",
			"j*ck",
			false,
		},
		{
			"Colon as a username is invalid because we use colons as delimiters in DB keys",
			":",
			false,
		},
		{
			"Colon in the middle of a username is invalid",
			"am:pm",
			false,
		},
		{
			"handle with leading @ is invalid",
			"@jack",
			false,
		},
		{
			"handle with more than 60 characters is invalid",
			strings.Repeat("A", 61),
			false,
		},
		{
			"handle with illegal characters is invalid",
			"jack.and.jill",
			false,
		},
	}

	for _, tt := range tests {
		validActual := Username(tt.username)
		if validActual != tt.validExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.username, validActual, tt.validExpected)
		}
	}
}
