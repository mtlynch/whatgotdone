package types

import "time"

// JournalEntry represents a user's What Got Done update. The entry can be
// public or a private draft that has not yet been published.
type JournalEntry struct {
	Date             string    `json:"date" yaml:"date" firestore:"date,omitempty"`
	CreationTime     time.Time `json:"creationTime" yaml:"creationTime" firestore:"creationTime"`
	LastModifiedTime time.Time `json:"lastModifiedTime" yaml:"lastModifiedTime" firestore:"lastModifiedTime"`
	Markdown         string    `json:"markdown" yaml:"markdown" firestore:"markdown,omitempty"`
}
