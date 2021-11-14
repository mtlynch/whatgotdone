package firestore

import (
	"errors"
	"log"

	"google.golang.org/api/iterator"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetEntry returns the published entry for the given date.
func (c client) GetEntry(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	iter := c.firestoreClient.Collection(entriesRootKey).Doc(string(username)).Collection(perUserEntriesKey).Documents(c.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return types.JournalEntry{}, err
		}
		var j types.JournalEntry
		doc.DataTo(&j)
		if j.Date == date {
			return j, nil
		}
	}
	return types.JournalEntry{}, datastore.EntryNotFoundError{
		Username: username,
		Date:     date,
	}
}

func (c client) ReadEntries(datastore.EntryFilter) ([]types.JournalEntry, error) {
	return []types.JournalEntry{}, errors.New("not implemented")
}

// InsertEntry saves an entry to the datastore, overwriting any existing entry
// with the same name and username.
func (c client) InsertEntry(username types.Username, j types.JournalEntry) error {
	log.Printf("adding new entry for %s: %+v", username, j)
	// Create a User document so that its children appear in Firestore console.
	c.firestoreClient.Collection(entriesRootKey).Doc(string(username)).Set(c.ctx, userDocument{
		Username:     username,
		LastModified: j.LastModified,
	})
	_, err := c.firestoreClient.Collection(entriesRootKey).Doc(string(username)).Collection(perUserEntriesKey).Doc(string(j.Date)).Set(c.ctx, j)
	return err
}
