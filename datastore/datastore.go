package datastore

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/mtlynch/whatgotdone/types"
)

type Datastore interface {
	Users() ([]string, error)
	All(username string) ([]types.JournalEntry, error)
	GetDraft(username string, date string) (types.JournalEntry, error)
	Insert(username string, j types.JournalEntry) error
	InsertDraft(username string, j types.JournalEntry) error
	GetReactions(entryAuthor string, entryDate string) ([]types.Reaction, error)
	AddReaction(entryAuthor string, entryDate string, reaction types.Reaction) error
	Close() error
}

type DraftNotFoundError struct {
	Username string
	Date     string
}

func (f DraftNotFoundError) Error() string {
	return fmt.Sprintf("Could not find draft entry for user %s on date %s", f.Username, f.Date)
}

type (
	defaultClient struct {
		firestoreClient *firestore.Client
		ctx             context.Context
	}

	userDocument struct {
		Username     string `firestore:"username,omitempty"`
		LastModified string `firestore:"lastModified,omitempty"`
	}

	reactionsDocument struct {
		Reactions []types.Reaction `firestore:"reactions,omitempty"`
	}

	entryReactionsDocument struct {
		entryAuthor string `firestore:"entryAuthor,omitempty"`
		entryDate   string `firestore:"entryDate,omitempty"`
	}
)

const (
	entriesRootKey      = "journalEntries"
	perUserEntriesKey   = "entries"
	draftsRootKey       = "journalDrafts"
	perUserDraftsKey    = "drafts"
	reactionsRootKey    = "reactions"
	perUserReactionsKey = "perUserReactions"
	secretsRootKey      = "secrets"
	secretUserKitDocKey = "userKitKey"
)

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
	iter := c.firestoreClient.Collection(entriesRootKey).Documents(c.ctx)
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

func (c defaultClient) All(username string) ([]types.JournalEntry, error) {
	entries := make([]types.JournalEntry, 0)
	iter := c.firestoreClient.Collection(entriesRootKey).Doc(username).Collection(perUserEntriesKey).Documents(c.ctx)
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
		if strings.TrimSpace(j.Markdown) == "" {
			continue
		}
		entries = append(entries, j)
	}
	return entries, nil
}

func (c defaultClient) GetDraft(username string, date string) (types.JournalEntry, error) {
	iter := c.firestoreClient.Collection(draftsRootKey).Doc(username).Collection(perUserDraftsKey).Documents(c.ctx)
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
	return types.JournalEntry{}, DraftNotFoundError{
		Username: username,
		Date:     date,
	}
}

func (c defaultClient) Insert(username string, j types.JournalEntry) error {
	// Create a User document so that its children appear in Firestore console.
	c.firestoreClient.Collection(entriesRootKey).Doc(username).Set(c.ctx, userDocument{
		Username:     username,
		LastModified: j.LastModified,
	})
	_, err := c.firestoreClient.Collection(entriesRootKey).Doc(username).Collection(perUserEntriesKey).Doc(j.Date).Set(c.ctx, j)
	return err
}

func (c defaultClient) Close() error {
	return c.firestoreClient.Close()
}

func (c defaultClient) InsertDraft(username string, j types.JournalEntry) error {
	// Create a User document so that its children appear in Firestore console.
	c.firestoreClient.Collection(draftsRootKey).Doc(username).Set(c.ctx, userDocument{
		Username:     username,
		LastModified: j.LastModified,
	})
	_, err := c.firestoreClient.Collection(draftsRootKey).Doc(username).Collection(perUserDraftsKey).Doc(j.Date).Set(c.ctx, j)
	return err
}

func (c defaultClient) GetReactions(entryAuthor string, entryDate string) ([]types.Reaction, error) {
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

func (c defaultClient) AddReaction(entryAuthor string, entryDate string, reaction types.Reaction) error {
	// Create a entryReactionsDocument document so that its children appear in Firestore console.
	c.firestoreClient.Collection(reactionsRootKey).Doc(getEntryReactionsKey(entryAuthor, entryDate)).Set(c.ctx, entryReactionsDocument{
		entryAuthor: entryAuthor,
		entryDate:   entryDate,
	})

	_, err := c.firestoreClient.Collection(reactionsRootKey).Doc(getEntryReactionsKey(entryAuthor, entryDate)).Collection(perUserReactionsKey).Doc(reaction.Username).Set(c.ctx, reaction)

	return err
}

func getEntryReactionsKey(entryAuthor string, entryDate string) string {
	return entryAuthor + ":" + entryDate
}
