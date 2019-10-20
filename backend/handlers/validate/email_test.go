package validate

import (
	"testing"
)

func TestEmailAddress(t *testing.T) {
	var tests = []struct {
		explanation   string
		email         string
		validExpected bool
	}{
		{
			"well-formed email address is valid",
			"hello@example.com",
			true,
		},
		{
			"empty string is invalid",
			"",
			false,
		},
		{
			"email without @ is invalid",
			"hello[at]example.com",
			false,
		},
		{
			"email with extraneous data is invalid",
			"Barry Gibbs <bg@example.com>",
			false,
		},
		{
			"email with angle brackets is invalid",
			"<hello@example.com>",
			false,
		},
	}

	for _, tt := range tests {
		validActual := EmailAddress(tt.email)
		if validActual != tt.validExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.email, validActual, tt.validExpected)
		}
	}
}
