package types

type (
	// Date in YYYY-MM-DD format.
	EntryDate string

	// JournalEntry represents a user's What Got Done update. The entry can be
	// public or a private draft that has not yet been published.
	JournalEntry struct {
		Author       Username
		Date         EntryDate `json:"date" yaml:"date" firestore:"date,omitempty"`
		LastModified string    `json:"lastModified" yaml:"lastModified" firestore:"lastModified,omitempty"`
		Markdown     string    `json:"markdown" yaml:"markdown" firestore:"markdown,omitempty"`
	}
)
