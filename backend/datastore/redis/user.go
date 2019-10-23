package redis

import (
	"errors"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// Users returns all the users who have published entries.
func (c client) Users() (users []string, err error) {
	return []string{}, errors.New("not implemented")
}

// UserProfile returns profile information about the given user.
func (c client) GetUserProfile(username string) (types.UserProfile, error) {
	return types.UserProfile{}, errors.New("not implemented")
}

// SetUserProfile updates the given user's profile or creates a new profile for
// the user.
func (c client) SetUserProfile(username string, p types.UserProfile) error {
	return errors.New("not implemented")
}
