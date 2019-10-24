package validate

import (
	"regexp"
)

// Username validates that a What Got Done username is valid.
// Valid What Got Done usernames are whatever UserKit allows, which is
// currently:
//  * 1-60 characters
//  * Only English letters, numbers, and underscores
func Username(username string) bool {
	// If something goes wrong in a JavaScript-based client, it will send the
	// literal string "undefined" as the username when the username variable is
	// undefined.
	if username == "undefined" {
		return false
	}

	re := regexp.MustCompile("^[A-Za-z0-9_]{1,60}$")
	return re.MatchString(username)
}
