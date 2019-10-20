package validate

import "regexp"

// UserBio validates that a user's profile bio is valid.
func UserBio(bio string) bool {
	invalidPatterns := []string{
		".{301}",    // excessive length
		"```",       // fenced code block
		"!\\[.*\\]", // image
		"(?m)^#",    // heading
	}
	for _, p := range invalidPatterns {
		if regexp.MustCompile(p).MatchString(bio) {
			return false
		}
	}

	return true
}
