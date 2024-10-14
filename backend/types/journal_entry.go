package types

import "time"

type (
	// Date in YYYY-MM-DD format.
	EntryDate string

	EntryContent string

	// JournalEntry represents a user's What Got Done update. The entry can be
	// public or a private draft that has not yet been published.
	JournalEntry struct {
		Author       Username
		Date         EntryDate    `json:"date"`
		LastModified time.Time    `json:"lastModified"`
		Markdown     EntryContent `json:"markdown"`
	}
)
