package types

type JournalEntry struct {
	Date         string `json:"date" firestore:"date,omitempty"`
	LastModified string `json:"lastModified" firestore:"lastModified,omitempty"`
	Markdown     string `json:"markdown" firestore:"markdown,omitempty"`
}
