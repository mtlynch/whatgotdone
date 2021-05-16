// Package firestore implements a datastore.Datastore interface using Google
// Cloud Firestore as a backend.
package firestore

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	fc "github.com/mtlynch/whatgotdone/backend/datastore/firestore/client"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type (
	client struct {
		firestoreClient *firestore.Client
		ctx             context.Context
	}

	userDocument struct {
		Username     types.Username `firestore:"username,omitempty"`
		LastModified string         `firestore:"lastModified,omitempty"`
	}

	reactionsDocument struct {
		Reactions []types.Reaction `firestore:"reactions,omitempty"`
	}

	pageViewsDocument struct {
		Path  string `firestore:"path"`
		Views int    `firestore:"views"`
	}

	entryReactionsDocument struct {
		entryAuthor types.Username `firestore:"entryAuthor,omitempty"`
		entryDate   string         `firestore:"entryDate,omitempty"`
	}

	followDocument struct {
		Follower     types.Username `firestore:"follower"`
		LastModified time.Time      `firestore:"lastModified"`
	}
)

const (
	entriesRootKey      = "journalEntries"
	perUserEntriesKey   = "entries"
	draftsRootKey       = "journalDrafts"
	perUserDraftsKey    = "drafts"
	pageViewsRootKey    = "pageViews"
	preferencesRootKey  = "preferences"
	reactionsRootKey    = "reactions"
	perUserReactionsKey = "perUserReactions"
	secretsRootKey      = "secrets"
	secretUserKitDocKey = "userKitKey"
	userProfilesRootKey = "userProfiles"
	followingRootKey    = "following"
	perUserFollowingKey = "perUserFollowers"
)

// New creates a new Datastore instance.
func New() datastore.Datastore {
	ctx := context.Background()

	c, err := fc.New(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return &client{
		firestoreClient: c,
		ctx:             ctx,
	}
}
