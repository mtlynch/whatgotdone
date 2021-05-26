package parse

import (
	"testing"

	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestEmailAddress(t *testing.T) {
	var tests = []struct {
		explanation    string
		email          string
		validExpected  bool
		parsedExpected types.EmailAddress
	}{
		{
			"well-formed email address is valid",
			"hello@example.com",
			true,
			"hello@example.com",
		},
		{
			"email with name data is valid",
			"Barry Gibbs <bg@example.com>",
			true,
			"bg@example.com",
		},
		{
			"email with angle brackets is valid",
			"<hello@example.com>",
			true,
			"hello@example.com",
		},
		{
			"empty string is invalid",
			"",
			false,
			"",
		},
		{
			"email without @ is invalid",
			"hello[at]example.com",
			false,
			"",
		},
	}

	for _, tt := range tests {
		parsedActual, errActual := EmailAddress(tt.email)
		if (errActual == nil) != tt.validExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.email, errActual, tt.validExpected)
		} else if parsedActual != tt.parsedExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.email, parsedActual, tt.parsedExpected)
		}
	}
}
