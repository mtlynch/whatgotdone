// +build dev staging

package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/mtlynch/whatgotdone/backend/datastore/firestore/client"
)

type wiper struct {
	firestoreClient *firestore.Client
	ctx             context.Context
}

func newWiper(ctx context.Context) wiper {
	c, err := client.New(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return wiper{
		firestoreClient: c,
		ctx:             ctx,
	}
}

const (
	entriesRootKey      = "journalEntries"
	draftsRootKey       = "journalDrafts"
	pageViewsRootKey    = "pageViews"
	preferencesRootKey  = "preferences"
	reactionsRootKey    = "reactions"
	secretsRootKey      = "secrets"
	userProfilesRootKey = "userProfiles"
	followingRootKey    = "following"
)

func (w *wiper) Wipe() {
	rootKeys := []string{
		entriesRootKey,
		draftsRootKey,
		pageViewsRootKey,
		preferencesRootKey,
		reactionsRootKey,
		secretsRootKey,
		userProfilesRootKey,
		followingRootKey,
	}
	for _, collectionKey := range rootKeys {
		if err := w.deleteCollection(collectionKey); err != nil {
			log.Printf("failed to delete %s", collectionKey)
		}
	}
}

func (w *wiper) deleteCollection(collectionKey string) error {
	ref := w.firestoreClient.Collection(collectionKey)
	for {
		iter := ref.Limit(50).Documents(w.ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := w.firestoreClient.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(w.ctx)
		if err != nil {
			return err
		}
	}
}
