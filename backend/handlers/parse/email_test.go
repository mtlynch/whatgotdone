package parse

import (
	"fmt"
	"testing"

	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestEmailAddress(t *testing.T) {
	for _, tt := range []struct {
		explanation    string
		input          string
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
	} {
		t.Run(fmt.Sprintf("%s: %s", tt.explanation, tt.input), func(t *testing.T) {
			parsed, err := EmailAddress(tt.input)
			if got, want := (err == nil), tt.validExpected; got != want {
				t.Fatalf("valid=%v, want=%v", got, want)
			}
			if got, want := parsed, tt.parsedExpected; got != want {
				t.Errorf("email=%v, want=%v", got, want)
			}
		})
	}
}
