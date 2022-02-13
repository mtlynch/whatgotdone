package parse

import (
	"errors"

	"github.com/mtlynch/social-go/v2/social"

	"github.com/mtlynch/whatgotdone/backend/types"
)

var (
	ErrInvalidTwitterHandle = errors.New("invalid twitter handle")
)

// TwitterHandle parses a raw string into a TwitterHandle, validating that it's
// a legal Twitter handle.
func TwitterHandle(handle string) (types.TwitterHandle, error) {
	if isUndefinedFromJavascript(handle) {
		return types.TwitterHandle(""), ErrInvalidTwitterHandle
	}

	h, err := social.ParseTwitterHandle(handle)
	if err != nil {
		return types.TwitterHandle(""), ErrInvalidTwitterHandle
	}

	return types.TwitterHandle(h), nil
}
