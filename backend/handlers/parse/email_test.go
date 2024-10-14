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
		parsedExpected types.EmailAddress
		errText        string
	}{
		{
			"well-formed email address is valid",
			"hello@example.com",
			"hello@example.com",
			"",
		},
		{
			"email with name data is valid",
			"Barry Gibbs <bg@example.com>",
			"bg@example.com",
			"",
		},
		{
			"email with angle brackets is valid",
			"<hello@example.com>",
			"hello@example.com",
			"",
		},
		{
			"empty string is invalid",
			"",
			"",
			"email address cannot be empty",
		},
		{
			"email without @ is invalid",
			"hello[at]example.com",
			"",
			"email address must have @ symbol",
		},
	} {
		t.Run(fmt.Sprintf("%s: %s", tt.explanation, tt.input), func(t *testing.T) {
			parsed, err := EmailAddress(tt.input)
			if got, want := errToString(err), tt.errText; got != want {
				t.Fatalf("err=%v, want=%v", got, want)
			}
			if got, want := parsed, tt.parsedExpected; got != want {
				t.Errorf("email=%v, want=%v", got, want)
			}
		})
	}
}

func errToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
