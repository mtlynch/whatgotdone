package firestore

import (
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// Users returns all the users who have published entries.
func (c client) Users() (users []string, err error) {
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

// UserProfile returns profile information about the given user.
func (c client) GetUserProfile(username string) (profile types.UserProfile, err error) {
	doc := c.firestoreClient.Collection(userProfilesRootKey).Doc(username)
	docsnap, err := doc.Get(c.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return types.UserProfile{}, datastore.UserProfileNotFoundError{Username: username}
		}
		return types.UserProfile{}, err
	}
	var p types.UserProfile
	if err := docsnap.DataTo(&p); err != nil {
		return types.UserProfile{}, err
	}
	return p, nil
}

// SetUserProfile updates the given user's profile or creates a new profile for
// the user.
func (c client) SetUserProfile(username string, p types.UserProfile) error {
	_, err := c.firestoreClient.Collection(userProfilesRootKey).Doc(username).Set(c.ctx, p)
	return err
}
