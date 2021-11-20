package firestore

import (
	"log"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// Users returns all the users who have saved drafts or published profiles.
func (c client) Users() ([]types.Username, error) {
	usersWithDrafts, err := c.usernamesInCollection(draftsRootKey)
	if err != nil {
		return []types.Username{}, err
	}
	usersWithProfiles, err := c.usernamesInCollection(userProfilesRootKey)
	if err != nil {
		return []types.Username{}, err
	}
	usersWithPreferences, err := c.usernamesInCollection(preferencesRootKey)
	if err != nil {
		return []types.Username{}, err
	}
	allUsers := map[types.Username]bool{}
	for _, u := range append(usersWithDrafts, append(usersWithProfiles, usersWithPreferences...)...) {
		allUsers[u] = true
	}
	uniqueUsers := make([]types.Username, len(allUsers))
	i := 0
	for u := range allUsers {
		uniqueUsers[i] = u
		i++
	}
	return uniqueUsers, nil
}

func (c client) usernamesInCollection(collectionKey string) ([]types.Username, error) {
	users := []types.Username{}
	iter := c.firestoreClient.Collection(collectionKey).Documents(c.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		users = append(users, types.Username(doc.Ref.ID))
	}
	return users, nil
}

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
