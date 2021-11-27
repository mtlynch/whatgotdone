package types

type (
	UserBio         string
	EmailAddress    string
	TwitterHandle   string
	MastodonAddress string

	// UserProfile represents public information about a What Got Done user.
	UserProfile struct {
		AboutMarkdown   UserBio         `json:"aboutMarkdown"`
		EmailAddress    EmailAddress    `json:"emailAddress"`
		TwitterHandle   TwitterHandle   `json:"twitterHandle"`
		MastodonAddress MastodonAddress `json:"mastodonAddress"`
	}
)
