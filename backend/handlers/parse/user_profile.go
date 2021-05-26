package parse

import (
	"strings"

	"github.com/mtlynch/whatgotdone/backend/types"
	"github.com/mtlynch/whatgotdone/backend/types/requests"
)

// ProfileUpdateRequest parses a raw profile update request into a UserProfile,
// validating that the request is well-formed and all fields are legal.
func ProfileUpdateRequest(pur requests.ProfileUpdate) (types.UserProfile, error) {
	p := types.UserProfile{}
	if pur.AboutMarkdown != "" {
		bio, err := UserBio(strings.TrimSpace(pur.AboutMarkdown))
		if err != nil {
			return types.UserProfile{}, err
		}
		p.AboutMarkdown = bio
	}
	if pur.EmailAddress != "" {
		email, err := EmailAddress(pur.EmailAddress)
		if err != nil {
			return types.UserProfile{}, err
		}
		p.EmailAddress = email
	}
	if pur.TwitterHandle != "" {
		twitterHandle, err := TwitterHandle(pur.TwitterHandle)
		if err != nil {
			return types.UserProfile{}, err
		}
		p.TwitterHandle = twitterHandle
	}
	return p, nil
}
