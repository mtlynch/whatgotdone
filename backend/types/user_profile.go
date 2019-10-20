package types

// UserProfile represents public information about a What Got Done user.
type UserProfile struct {
	AboutMarkdown string `json:"aboutMarkdown" firestore:"aboutMarkdown,omitempty"`
	EmailAddress  string `json:"emailAddress" firestore:"emailAddress,omitempty"`
	TwitterHandle string `json:"twitterHandle" firestore:"twitterHandle,omitempty"`
}
