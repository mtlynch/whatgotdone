package validate

import "regexp"

// TwitterHandle validates that a Twitter handle is valid.
func TwitterHandle(handle string) bool {
	re := regexp.MustCompile("^[A-Za-z0-9_]{1,15}$")
	return re.MatchString(handle)
}
