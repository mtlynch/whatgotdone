package parse

import (
	"errors"
	"regexp"

	"github.com/mtlynch/whatgotdone/backend/types"
)

var (
	ErrInvalidTwitterHandle = errors.New("invalid twitter handle")

	twitterHandlePattern = regexp.MustCompile("^[A-Za-z0-9_]{4,15}$")
)

// TwitterHandle parses a raw string into a TwitterHandle, validating that it's
// a legal Twitter handle.
func TwitterHandle(handle string) (types.TwitterHandle, error) {
	if isUndefinedFromJavascript(handle) {
		return types.TwitterHandle(""), ErrInvalidTwitterHandle
	}

	if !twitterHandlePattern.MatchString(handle) {
		return types.TwitterHandle(""), ErrInvalidTwitterHandle
	}
	return types.TwitterHandle(handle), nil
}
