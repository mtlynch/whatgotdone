package parse

import (
	"strings"

	"github.com/mtlynch/whatgotdone/backend/types"
	"github.com/mtlynch/whatgotdone/backend/types/requests"
)

func ProfileUpdateRequest(pur requests.ProfileUpdate) (types.UserProfile, error) {
	up := types.UserProfile{}
	if pur.AboutMarkdown != "" {
		bio, err := UserBio(strings.TrimSpace(pur.AboutMarkdown))
		if err != nil {
			return types.UserProfile{}, err
		}
		up.AboutMarkdown = bio
	}
	if pur.EmailAddress != "" {
		email, err := EmailAddress(pur.EmailAddress)
		if err != nil {
			return types.UserProfile{}, err
		}
		up.EmailAddress = email
	}
	if pur.TwitterHandle != "" {
		twitterHandle, err := TwitterHandle(pur.TwitterHandle)
		if err != nil {
			return types.UserProfile{}, err
		}
		up.TwitterHandle = twitterHandle
	}
	return up, nil
}
