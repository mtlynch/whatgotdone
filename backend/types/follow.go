package types

import "time"

// Follow represents a following relationship between a leader and a follower.
type Follow struct {
	Leader   string    `firestore:"leader"`
	Follower string    `firestore:"follower"`
	Created  time.Time `firestore:"created"`
}
