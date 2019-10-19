package firestore

import (
	"google.golang.org/api/iterator"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetDraft returns an entry draft for the given user for the given date.
func (c client) GetDraft(username string, date string) (types.JournalEntry, error) {
	iter := c.firestoreClient.Collection(draftsRootKey).Doc(username).Collection(perUserDraftsKey).Documents(c.ctx)
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
	return types.JournalEntry{}, datastore.DraftNotFoundError{
		Username: username,
		Date:     date,
	}
}

// InsertDraft saves an entry draft to the datastore, overwriting any existing
// entry with the same name and username.
func (c client) InsertDraft(username string, j types.JournalEntry) error {
	// Create a User document so that its children appear in Firestore console.
	c.firestoreClient.Collection(draftsRootKey).Doc(username).Set(c.ctx, userDocument{
		Username:     username,
		LastModified: j.LastModified,
	})
	_, err := c.firestoreClient.Collection(draftsRootKey).Doc(username).Collection(perUserDraftsKey).Doc(j.Date).Set(c.ctx, j)
	return err
}
