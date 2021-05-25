package parse

import (
	"errors"
	"regexp"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// TwitterHandle parses a raw string into a TwitterHandle, validating that it's
// a legal Twitter handle.
func TwitterHandle(handle string) (types.TwitterHandle, error) {
	if !regexp.MustCompile("^[A-Za-z0-9_]{1,15}$").MatchString(handle) {
		return types.TwitterHandle(""), errors.New("invalid Twitter handle")
	}
	return types.TwitterHandle(handle), nil
}
