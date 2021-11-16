package firestore

import (
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// UserProfile returns profile information about the given user.
func (c client) GetUserProfile(username types.Username) (types.UserProfile, error) {
	doc := c.firestoreClient.Collection(userProfilesRootKey).Doc(string(username))
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
func (c client) SetUserProfile(username types.Username, p types.UserProfile) error {
	log.Printf("saving user profile to datastore: %s -> %+v", username, p)
	_, err := c.firestoreClient.Collection(userProfilesRootKey).Doc(string(username)).Set(c.ctx, p)
	return err
}
