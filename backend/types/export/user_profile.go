package export

import "github.com/mtlynch/whatgotdone/backend/types"

type UserProfile struct {
	AboutMarkdown   types.UserBio         `json:"aboutMarkdown" firestore:"about_markdown"`
	EmailAddress    types.EmailAddress    `json:"emailAddress" firestore:"email_address"`
	TwitterHandle   types.TwitterHandle   `json:"twitterHandle" firestore:"twitter_handle"`
	MastodonAddress types.MastodonAddress `json:"mastodonAddress" yaml:"mastodon_address"`
}
