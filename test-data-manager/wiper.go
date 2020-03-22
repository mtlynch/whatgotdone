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
	w.deleteAllCollections()
}

func (w *wiper) deleteAllCollections() error {
	iter := w.firestoreClient.Collections(w.ctx)
	for {
		collection, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		err = w.deleteCollection(collection)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *wiper) deleteCollection(collection *firestore.CollectionRef) error {
	for {
		iter := collection.Limit(50).Documents(w.ctx)
		numDeleted := 0

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			err = w.recursiveDeleteDocument(doc.Ref)
			if err != nil {
				return err
			}
			numDeleted++
		}

		// If there are no documents to delete, the process is over.
		if numDeleted == 0 {
			return nil
		}
	}
}

func (w *wiper) recursiveDeleteDocument(doc *firestore.DocumentRef) error {
	iter := doc.Collections(w.ctx)
	for {
		collection, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		err = w.deleteCollection(collection)
		if err != nil {
			return err
		}
	}

	_, err := doc.Delete(w.ctx)
	return err
}
