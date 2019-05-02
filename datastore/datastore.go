package datastore

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/mtlynch/whatgotdone/types"
)

type Datastore interface {
	All(username string) ([]types.JournalEntry, error)
	Get(username string, date string) (types.JournalEntry, error)
	Insert(username string, j types.JournalEntry) error
	Close() error
}

type EntryNotFoundError struct {
	Username string
	Date     string
}

func (f EntryNotFoundError) Error() string {
	return fmt.Sprintf("Could not find journal entry for user %s on date %s", f.Username, f.Date)
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

func (c defaultClient) All(username string) (entries []types.JournalEntry, err error) {
	iter := c.firestoreClient.Collection("journalEntries").Doc(username).Collection("entries").Documents(c.ctx)
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

func (c defaultClient) Get(username string, date string) (types.JournalEntry, error) {
	iter := c.firestoreClient.Collection("journalEntries").Doc(username).Collection("entries").Documents(c.ctx)
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
	return types.JournalEntry{}, EntryNotFoundError{
		Username: username,
		Date:     date,
	}
}

func (c defaultClient) Insert(username string, j types.JournalEntry) error {
	_, err := c.firestoreClient.Collection("journalEntries").Doc(username).Collection("entries").Doc(j.Date).Set(c.ctx, j)
	return err
}

func (c defaultClient) Close() error {
	return c.firestoreClient.Close()
}
