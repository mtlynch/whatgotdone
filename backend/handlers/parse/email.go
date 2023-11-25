package parse

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/mtlynch/whatgotdone/backend/types"
)

var (
	ErrEmailEmpty           = errors.New("email address cannot be empty")
	ErrEmailMissingAtSymbol = errors.New("email address must have @ symbol")
)

// EmailAddress parses a raw email string into an EmailAddress object,
// validating that it's well-formed.
func EmailAddress(email string) (types.EmailAddress, error) {
	// mail.ParseAddress will catch these issues as well, but we catch them
	// early at this layer for clearer error messages.
	if strings.TrimSpace(email) == "" {
		return types.EmailAddress(""), ErrEmailEmpty
	}
	if !strings.ContainsAny(email, "@") {
		return types.EmailAddress(""), ErrEmailMissingAtSymbol
	}

	a, err := mail.ParseAddress(email)
	if err != nil {
		return types.EmailAddress(""), err
	}
	return types.EmailAddress(a.Address), nil
}
