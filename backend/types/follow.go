package types

import "time"

// Follow represents a following relationship between a leader and a follower.
type Follow struct {
	Leader   Username  `firestore:"leader"`
	Follower Username  `firestore:"follower"`
	Created  time.Time `firestore:"created"`
}
