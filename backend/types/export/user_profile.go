package export

import "github.com/mtlynch/whatgotdone/backend/types"

type UserProfile struct {
	AboutMarkdown   types.UserBio         `json:"aboutMarkdown" yaml:"about_markdown"`
	EmailAddress    types.EmailAddress    `json:"emailAddress" yaml:"email_address"`
	TwitterHandle   types.TwitterHandle   `json:"twitterHandle" yaml:"twitter_handle"`
	MastodonAddress types.MastodonAddress `json:"mastodonAddress" yaml:"mastodon_address"`
}
