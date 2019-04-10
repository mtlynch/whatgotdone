package datastore

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/mtlynch/whatgotdone/types"
)

type Datastore interface {
	All() ([]types.JournalEntry, error)
	InsertJournalEntry(types.JournalEntry) error
	Close() error
}

type defaultClient struct {
	firestoreClient *firestore.Client
	ctx             context.Context
}

const devServiceAccount = "service-account-creds.json"
const firestoreProjectId = "whatgotdone"

func newFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	if _, err := os.Stat(devServiceAccount); !os.IsNotExist(err) {
		return firestore.NewClient(ctx, firestoreProjectId, option.WithCredentialsFile(devServiceAccount))
	}
	return firestore.NewClient(ctx, firestoreProjectId)
}

func New() Datastore {
	ctx := context.Background()

	client, err := newFirestoreClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return &defaultClient{
		firestoreClient: client,
		ctx:             ctx,
	}
}

func (c defaultClient) All() (entries []types.JournalEntry, err error) {
	iter := c.firestoreClient.Collection("journalEntries").Documents(c.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var j types.JournalEntry
		doc.DataTo(&j)
		entries = append(entries, j)
	}
	return entries, nil
}

func (c defaultClient) InsertJournalEntry(j types.JournalEntry) error {
	_, err := c.firestoreClient.Collection("journalEntries").Doc(j.Date).Set(c.ctx, j)
	return err
}

func (c defaultClient) Close() error {
	return c.firestoreClient.Close()
}
