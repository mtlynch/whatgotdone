package validate

import (
	"testing"
)

func TestEmail(t *testing.T) {
	var tests = []struct {
		explanation   string
		email          string
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
	}

	for _, tt := range tests {
		validActual := Email(tt.date)
		if validActual != tt.validExpected {
			t.Errorf("%s: input: %s, got %v, want %v", tt.explanation, tt.email, validActual, tt.validExpected)
		}
	}
}