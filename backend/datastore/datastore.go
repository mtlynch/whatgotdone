package datastore

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/mtlynch/whatgotdone/backend/types"
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

func New() Datastore {
	ctx := context.Background()
	var client *firestore.Client
	var err error
	if _, err := os.Stat(devServiceAccount); !os.IsNotExist(err) {
		opt := option.WithCredentialsFile(devServiceAccount)
		client, err = firestore.NewClient(ctx, "whatgotdone", opt)
	} else {
		client, err = firestore.NewClient(ctx, "whatgotdone")
	}
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
