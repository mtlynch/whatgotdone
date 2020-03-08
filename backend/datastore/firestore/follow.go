package firestore

import (
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// InsertFollow adds a following relationship to the datastore.
func (c client) InsertFollow(leader, follower string) error {
	// Create a followDocument so that its children appear in Firestore console.
	c.firestoreClient.Collection(followingRootKey).Doc(follower).Set(c.ctx, followDocument{
		Follower:     follower,
		LastModified: time.Now().UTC(),
	})
	f := types.Follow{
		Leader:   leader,
		Follower: follower,
		Created:  time.Now().UTC(),
	}
	log.Printf("inserting new follow: %+v", f)
	_, err := c.firestoreClient.Collection(followingRootKey).Doc(follower).Collection(perUserFollowingKey).Doc(leader).Create(c.ctx, f)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return datastore.FollowAlreadyExistsError{
				Leader:   leader,
				Follower: follower}
		}
	}
	return err
}

// DeleteFollow removes a following relationship from the datastore.
func (c client) DeleteFollow(leader, follower string) error {
	log.Printf("deleting follow: %s -> %s", follower, leader)
	_, err := c.firestoreClient.Collection(followingRootKey).Doc(follower).Collection(perUserFollowingKey).Doc(leader).Delete(c.ctx)
	return err
}

// Following returns all the users the specified user is following.
func (c client) Following(follower string) ([]string, error) {
	docs, err := c.firestoreClient.Collection(followingRootKey).Doc(follower).Collection(perUserFollowingKey).Documents(c.ctx).GetAll()
	if err != nil {
		return nil, err
	}
	following := make([]string, len(docs))
	for i, doc := range docs {
		var f types.Follow
		doc.DataTo(&f)
		following[i] = f.Leader
	}
	return following, nil
}
