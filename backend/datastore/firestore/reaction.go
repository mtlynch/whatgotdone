package firestore

import (
	"log"

	"google.golang.org/api/iterator"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetReactions retrieves reader reactions associated with a published entry.
func (c client) GetReactions(entryAuthor types.Username, entryDate types.EntryDate) ([]types.Reaction, error) {
	reactions := []types.Reaction{}
	iter := c.firestoreClient.Collection(reactionsRootKey).Doc(getEntryReactionsKey(entryAuthor, entryDate)).Collection(perUserReactionsKey).Documents(c.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var r types.Reaction
		doc.DataTo(&r)
		reactions = append(reactions, r)
	}
	return reactions, nil
}

// AddReaction saves a reader reaction associated with a published entry,
// overwriting any existing reaction.
func (c client) AddReaction(entryAuthor types.Username, entryDate types.EntryDate, reaction types.Reaction) error {
	// Create a entryReactionsDocument document so that its children appear in Firestore console.
	c.firestoreClient.Collection(reactionsRootKey).Doc(getEntryReactionsKey(entryAuthor, entryDate)).Set(c.ctx, entryReactionsDocument{
		entryAuthor: entryAuthor,
		entryDate:   entryDate,
	})
	key := getEntryReactionsKey(entryAuthor, entryDate)
	log.Printf("adding reaction to datastore: %s -> %+v", key, reaction)
	_, err := c.firestoreClient.Collection(reactionsRootKey).Doc(key).Collection(perUserReactionsKey).Doc(string(reaction.Username)).Set(c.ctx, reaction)

	return err
}

// DeleteReaction removes a user's reaction to a published entry.
func (c client) DeleteReaction(entryAuthor types.Username, entryDate types.EntryDate, reactingUser types.Username) error {
	key := getEntryReactionsKey(entryAuthor, entryDate)
	log.Printf("deleting %s's reaction to %s from datastore", reactingUser, key)
	_, err := c.firestoreClient.Collection(reactionsRootKey).Doc(key).Collection(perUserReactionsKey).Doc(string(reactingUser)).Set(c.ctx, types.Reaction{
		Username: reactingUser,
		Symbol:   "",
	})

	return err
}

func getEntryReactionsKey(entryAuthor types.Username, entryDate types.EntryDate) string {
	return string(entryAuthor) + ":" + string(entryDate)
}
