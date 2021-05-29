package requests

type ProfileUpdate struct {
	AboutMarkdown   string `json:"aboutMarkdown"`
	EmailAddress    string `json:"emailAddress"`
	TwitterHandle   string `json:"twitterHandle"`
	MastodonAddress string `json:"mastodonAddress"`
}
