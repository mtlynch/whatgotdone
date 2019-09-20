package types

// JournalEntry represents a user's What Got Done update. The entry can be
// public or a private draft that has not yet been published.
type JournalEntry struct {
	Date         string `json:"date" firestore:"date,omitempty"`
	LastModified string `json:"lastModified" firestore:"lastModified,omitempty"`
	Markdown     string `json:"markdown" firestore:"markdown,omitempty"`
}
