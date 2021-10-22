package parse

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// illegalMastodonAddressCharacters represents characters that we reject in a
// Mastodon address.
const illegalMastodonAddressCharacters = " <>"

// MastodonAddress parses a raw Mastodon address string into a MastodonAddress
// object validating that it's well-formed.
func MastodonAddress(address string) (types.MastodonAddress, error) {
	// A Mastodon address follows the same rules as an email address, except it
	// can't have a name portion separated by angle brackets.
	if strings.ContainsAny(address, illegalMastodonAddressCharacters) {
		return types.MastodonAddress(""), errors.New("Mastodon address contains illegal characters")
	}

	a, err := mail.ParseAddress(address)
	if err != nil {
		return types.MastodonAddress(""), errors.New("invalid Mastodon address, must be in the form of handle@hostname")
	}
	return types.MastodonAddress(a.Address), nil
}
