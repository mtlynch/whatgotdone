package firestore

import (
	"log"
	"time"

	"google.golang.org/api/iterator"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetReactions retrieves reader reactions associated with a published entry.
func (c client) GetReactions(entryAuthor string, entryDate string) ([]types.Reaction, error) {
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
func (c client) AddReaction(entryAuthor string, entryDate string, reaction types.Reaction) error {
	// Create a entryReactionsDocument document so that its children appear in Firestore console.
	c.firestoreClient.Collection(reactionsRootKey).Doc(getEntryReactionsKey(entryAuthor, entryDate)).Set(c.ctx, entryReactionsDocument{
		entryAuthor: entryAuthor,
		entryDate:   entryDate,
	})

	reaction.CreationTime = time.Now()
	key := getEntryReactionsKey(entryAuthor, entryDate)
	log.Printf("adding reaction to datastore: %s -> %+v", key, reaction)
	_, err := c.firestoreClient.Collection(reactionsRootKey).Doc(key).Collection(perUserReactionsKey).Doc(reaction.Username).Set(c.ctx, reaction)

	return err
}

func getEntryReactionsKey(entryAuthor string, entryDate string) string {
	return entryAuthor + ":" + entryDate
}
