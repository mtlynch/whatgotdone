package firestore

import (
	"log"

	"google.golang.org/api/iterator"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetReactions retrieves reader reactions associated with a published entry.
func (c client) GetReactions(entryAuthor types.Username, entryDate string) ([]types.Reaction, error) {
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
func (c client) AddReaction(entryAuthor types.Username, entryDate string, reaction types.Reaction) error {
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

func getEntryReactionsKey(entryAuthor types.Username, entryDate string) string {
	return string(entryAuthor) + ":" + entryDate
}
