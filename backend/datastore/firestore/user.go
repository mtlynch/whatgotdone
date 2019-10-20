package firestore

import (
	"google.golang.org/api/iterator"

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
		return profile, err
	}
	var p types.UserProfile
	if err := docsnap.DataTo(&p); err != nil {
		return profile, err
	}
	return p, nil
}

// SetUserProfile updates the given user's profile.
func (c client) SetUserProfile(username string, p types.UserProfile) error {
	_, err := c.firestoreClient.Collection(userProfilesRootKey).Doc(username).Set(c.ctx, p)
	return err
}
