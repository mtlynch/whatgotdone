// Package datastore provides functionality for storing and retrieving
// persistent data. This package does not enforce access control, so it is the
// client's responsibility to enforce authentication and authorization before
// retrieving private data from the datastore.
package datastore

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// Datastore represents the Firestore datastore. It's responsible for storing
// and retrieving all persistent data (journal entries, journal drafts,
// reactions).
type Datastore interface {
	Users() ([]string, error)
	AllEntries(username string) ([]types.JournalEntry, error)
	GetDraft(username string, date string) (types.JournalEntry, error)
	InsertEntry(username string, j types.JournalEntry) error
	InsertDraft(username string, j types.JournalEntry) error
	GetReactions(entryAuthor string, entryDate string) ([]types.Reaction, error)
	AddReaction(entryAuthor string, entryDate string, reaction types.Reaction) error
	Close() error
}

// DraftNotFoundError occurs when no draft exists for a user with a given date.
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

// New creates a new Datastore instance.
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

// Users returns all the users who have published entries.
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

// AllEntries returns all published entries.
func (c defaultClient) AllEntries(username string) ([]types.JournalEntry, error) {
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

// GetDraft returns an entry draft for the given user for the given date.
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

// InsertEntry saves an entry to the datastore, overwriting any existing entry
// with the same name and username.
func (c defaultClient) InsertEntry(username string, j types.JournalEntry) error {
	// Create a User document so that its children appear in Firestore console.
	c.firestoreClient.Collection(entriesRootKey).Doc(username).Set(c.ctx, userDocument{
		Username:     username,
		LastModified: j.LastModified,
	})
	_, err := c.firestoreClient.Collection(entriesRootKey).Doc(username).Collection(perUserEntriesKey).Doc(j.Date).Set(c.ctx, j)
	return err
}

// Close cleans up datastore resources. Clients should not call any Datastore
// functions after calling Close().
func (c defaultClient) Close() error {
	return c.firestoreClient.Close()
}

// InsertDraft saves an entry draft to the datastore, overwriting any existing
// entry with the same name and username.
func (c defaultClient) InsertDraft(username string, j types.JournalEntry) error {
	// Create a User document so that its children appear in Firestore console.
	c.firestoreClient.Collection(draftsRootKey).Doc(username).Set(c.ctx, userDocument{
		Username:     username,
		LastModified: j.LastModified,
	})
	_, err := c.firestoreClient.Collection(draftsRootKey).Doc(username).Collection(perUserDraftsKey).Doc(j.Date).Set(c.ctx, j)
	return err
}

// GetReactions retrieves reader reactions associated with a published entry.
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

// AddReaction saves a reader reactions associated with a published entry,
// overwriting any existing reaction.
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

func getGoogleCloudProjectId() string {
	projectId := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectId == "" {
		log.Fatalf("GOOGLE_CLOUD_PROJECT environment variable must be set")
	}
	return projectId
}