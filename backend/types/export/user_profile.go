package export

import "github.com/mtlynch/whatgotdone/backend/types"

type UserProfile struct {
	AboutMarkdown   types.UserBio         `json:"aboutMarkdown"`
	EmailAddress    types.EmailAddress    `json:"emailAddress"`
	TwitterHandle   types.TwitterHandle   `json:"twitterHandle"`
	MastodonAddress types.MastodonAddress `json:"mastodonAddress"`
}
