package parse

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// UserBioMaxLength is the maximum allowable length of a user bio (in
// characters).
const UserBioMaxLength = 300

func UserBio(bio string) (types.UserBio, error) {
	invalidPatterns := []string{
		fmt.Sprintf(".{%d}", UserBioMaxLength+1), // excessive length
		"```",                                    // fenced code block
		"!\\[.*\\]",                              // image
		"(?m)^#",                                 // heading
	}
	for _, p := range invalidPatterns {
		if regexp.MustCompile(p).MatchString(bio) {
			return types.UserBio(""), errors.New("bio contains disallowed Markdown characters")
		}
	}
	return types.UserBio(bio), nil
}
