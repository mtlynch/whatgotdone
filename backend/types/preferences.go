package types

// Preferences are the set of a user's options for using the site.
type Preferences struct {
	Username      string `firestore:"username,omitempty"`
	EntryTemplate string `firestore:"entryTemplate,omitempty"`
}
