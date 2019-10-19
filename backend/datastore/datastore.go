// Package datastore provides functionality for storing and retrieving
// persistent data. This package does not enforce access control, so it is the
// client's responsibility to enforce authentication and authorization before
// retrieving private data from the datastore.
package datastore

import (
	"fmt"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// Datastore represents the What Got Done datastore. It's responsible for
// storing and retrieving all persistent data (journal entries, journal drafts,
// reactions).
type Datastore interface {
	Users() ([]string, error)
	AllEntries(username string) ([]types.JournalEntry, error)
	GetDraft(username string, date string) (types.JournalEntry, error)
	InsertEntry(username string, j types.JournalEntry) error
	InsertDraft(username string, j types.JournalEntry) error
	GetReactions(entryAuthor string, entryDate string) ([]types.Reaction, error)
	AddReaction(entryAuthor string, entryDate string, reaction types.Reaction) error
	Close() error
}

// DraftNotFoundError occurs when no draft exists for a user with a given date.
type DraftNotFoundError struct {
	Username string
	Date     string
}

func (f DraftNotFoundError) Error() string {
	return fmt.Sprintf("Could not find draft entry for user %s on date %s", f.Username, f.Date)
}
