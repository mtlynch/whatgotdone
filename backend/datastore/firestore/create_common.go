package firestore

import (
	"context"
	"log"

	"github.com/mtlynch/whatgotdone/backend/datastore"
)

// New creates a new Datastore instance.
func New(string) datastore.Datastore {
	ctx := context.Background()

	c, err := newFirestoreClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return &client{
		firestoreClient: c,
		ctx:             ctx,
	}
}
