package validate

import (
	"net/mail"
	"strings"
)

// illegalCharacters represents characters that are legal in an RFS 5322 email
// address, but that we do not allow.
const illegalCharacters = " <>"

// EmailAddress validates that an email address is well-formed.
func EmailAddress(email string) bool {
	if strings.ContainsAny(email, illegalCharacters) {
		return false
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}
