package sqlite

import "github.com/mtlynch/whatgotdone/backend/types"

// InsertFollow adds a following relationship to the datastore.
func (d db) InsertFollow(leader, follower types.Username) error {
	return notImplementedError
}

// DeleteFollow removes a following relationship from the datastore.
func (d db) DeleteFollow(leader, follower types.Username) error {
	return notImplementedError
}

// Followers returns all the users the specified user is following.
func (d db) Following(follower types.Username) ([]types.Username, error) {
	return []types.Username{}, notImplementedError
}
