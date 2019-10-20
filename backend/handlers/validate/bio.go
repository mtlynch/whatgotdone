package validate

import (
	"fmt"
	"regexp"
)

// UserBioMaxLength is the maximum allowable length of a user bio (in
// characters).
const UserBioMaxLength = 300

// UserBio validates that a user's profile bio is valid.
func UserBio(bio string) bool {
	invalidPatterns := []string{
		fmt.Sprintf(".{%d}", UserBioMaxLength+1), // excessive length
		"```",                                    // fenced code block
		"!\\[.*\\]",                              // image
		"(?m)^#",                                 // heading
	}
	for _, p := range invalidPatterns {
		if regexp.MustCompile(p).MatchString(bio) {
			return false
		}
	}

	return true
}
