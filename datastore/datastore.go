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
	Users() ([]string, error)
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

type (
	defaultClient struct {
		firestoreClient *firestore.Client
		ctx             context.Context
	}

	JournalEntry struct {
		Date         string `json:"date" firestore:"date,omitempty"`
		LastModified string `json:"lastModified" firestore:"lastModified,omitempty"`
		Markdown     string `json:"markdown" firestore:"markdown,omitempty"`
	}

	userDocument struct {
		Username     string `firestore:"username,omitempty"`
		LastModified string `firestore:"lastModified,omitempty"`
	}
)

func newFirestoreClient(ctx context.Context) (*firestore.Client, error) {

	const devServiceAccount = "service-account-creds.json"

	if _, err := os.Stat(devServiceAccount); !os.IsNotExist(err) {
		const firestoreProjectID = "whatgotdone-dev"
		return firestore.NewClient(ctx, firestoreProjectID, option.WithCredentialsFile(devServiceAccount))
	}
	const firestoreProjectID = "whatgotdone"
	return firestore.NewClient(ctx, firestoreProjectID)
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

func (c defaultClient) Users() (users []string, err error) {
	iter := c.firestoreClient.Collection("journalEntries").Documents(c.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var u userDocument
		doc.DataTo(&u)
		users = append(users, u.Username)
	}
	return users, nil
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
	// Create a User document so that its children appear in Firestore console.
	c.firestoreClient.Collection("journalEntries").Doc(username).Set(c.ctx, userDocument{
		Username:     username,
		LastModified: j.LastModified,
	})
	_, err := c.firestoreClient.Collection("journalEntries").Doc(username).Collection("entries").Doc(j.Date).Set(c.ctx, j)
	return err
}

func (c defaultClient) Close() error {
	return c.firestoreClient.Close()
}
