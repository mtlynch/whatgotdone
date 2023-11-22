package parse

import (
	"fmt"
	"testing"

	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestMastodonAddress(t *testing.T) {
	for _, tt := range []struct {
		explanation    string
		mastodon       string
		validExpected  bool
		parsedExpected types.MastodonAddress
	}{
		{
			"Mastodon address with leading @ is valid",
			"@hello@example.com",
			true,
			"hello@example.com",
		},
		{
			"well-formed Mastodon address is valid",
			"hello@example.com",
			true,
			"hello@example.com",
		},
		{
			"Mastodon with name data is invalid",
			"Barry Gibbs <bg@example.com>",
			false,
			"",
		},
		{
			"Mastodon with angle brackets is invalid",
			"<hello@example.com>",
			false,
			"",
		},
		{
			"empty string is invalid",
			"",
			false,
			"",
		},
		{
			"bare @ is invalid",
			"@",
			false,
			"",
		},
		{
			"double @ is invalid",
			"@@",
			false,
			"",
		},
		{
			"whitespace is invalid",
			"\t",
			false,
			"",
		},
		{
			"Mastodon address without @ is invalid",
			"hello[at]example.com",
			false,
			"",
		},
	} {
		t.Run(fmt.Sprintf("%s: %s", tt.explanation, tt.mastodon), func(t *testing.T) {
			parsed, err := MastodonAddress(tt.mastodon)
			if got, want := (err == nil), tt.validExpected; got != want {
				t.Fatalf("valid=%v, want=%v", got, want)
			}
			if got, want := parsed, tt.parsedExpected; got != want {
				t.Errorf("mastodon=%v, want=%v", got, want)
			}
		})
	}
}
