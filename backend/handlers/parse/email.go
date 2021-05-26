package parse

import (
	"net/mail"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// EmailAddress parses a raw email string into an EmailAddress object,
// validating that it's well-formed.
func EmailAddress(email string) (types.EmailAddress, error) {
	a, err := mail.ParseAddress(email)
	if err != nil {
		return types.EmailAddress(""), err
	}
	return types.EmailAddress(a.Address), nil
}
