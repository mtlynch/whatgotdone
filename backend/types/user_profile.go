package types

type (
	UserBio         string
	EmailAddress    string
	TwitterHandle   string
	MastodonAddress string

	// UserProfile represents public information about a What Got Done user.
	UserProfile struct {
		AboutMarkdown   UserBio         `json:"aboutMarkdown" firestore:"aboutMarkdown,omitempty"`
		EmailAddress    EmailAddress    `json:"emailAddress" firestore:"emailAddress,omitempty"`
		TwitterHandle   TwitterHandle   `json:"twitterHandle" firestore:"twitterHandle,omitempty"`
		MastodonAddress MastodonAddress `json:"mastodonAddress" firestore:"mastodonAddress,omitempty"`
	}
)
