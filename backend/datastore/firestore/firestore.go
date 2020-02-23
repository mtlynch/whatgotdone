// Package firestore implements a datastore.Datastore interface using Google
// Cloud Firestore as a backend.
package firestore

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/mtlynch/whatgotdone/backend/types"
)

type (
	client struct {
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

	pageViewsDocument struct {
		Path  string `firestore:"path"`
		Views int    `firestore:"views"`
	}

	entryReactionsDocument struct {
		entryAuthor string `firestore:"entryAuthor,omitempty"`
		entryDate   string `firestore:"entryDate,omitempty"`
	}

	followDocument struct {
		Follower     string    `firestore:"follower"`
		LastModified time.Time `firestore:"lastModified"`
	}
)

const (
	entriesRootKey      = "journalEntries"
	perUserEntriesKey   = "entries"
	draftsRootKey       = "journalDrafts"
	perUserDraftsKey    = "drafts"
	pageViewsRootKey    = "pageViews"
	reactionsRootKey    = "reactions"
	perUserReactionsKey = "perUserReactions"
	secretsRootKey      = "secrets"
	secretUserKitDocKey = "userKitKey"
	userProfilesRootKey = "userProfiles"
	followingRootKey    = "following"
	perUserFollowingKey = "perUserFollowers"
)

func getGoogleCloudProjectID() string {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatalf("GOOGLE_CLOUD_PROJECT environment variable must be set")
	}
	return projectID
}
