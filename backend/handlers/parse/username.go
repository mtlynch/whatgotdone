package parse

import (
	"errors"
	"regexp"

	"github.com/mtlynch/whatgotdone/backend/types"
)

var usernamePattern = regexp.MustCompile("^[A-Za-z0-9_]{1,60}$")

// Username parses a What Got Done username from a raw string.
// Valid What Got Done usernames are whatever UserKit allows, which is
// currently:
//  * 1-60 characters
//  * Only English letters, numbers, and underscores
func Username(username string) (types.Username, error) {
	// If something goes wrong in a JavaScript-based client, it will send the
	// literal string "undefined" as the username when the username variable is
	// undefined.
	if username == "undefined" {
		return "", errors.New("invalid username")
	}

	if !usernamePattern.MatchString(username) {
		return "", errors.New("invalid username")
	}
	return types.Username(username), nil
}
