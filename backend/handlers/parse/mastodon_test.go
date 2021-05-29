package parse

import (
	"testing"

	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestMastodonAddress(t *testing.T) {
	var tests = []struct {
		explanation    string
		mastodon       string
		validExpected  bool
		parsedExpected types.MastodonAddress
	}{
		{
			"well-formed Mastodon Address is valid",
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
	}

	for _, tt := range tests {
		parsedActual, errActual := MastodonAddress(tt.mastodon)
		if (errActual == nil) != tt.validExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.mastodon, errActual, tt.validExpected)
		} else if parsedActual != tt.parsedExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.mastodon, parsedActual, tt.parsedExpected)
		}
	}
}
