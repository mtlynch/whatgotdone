package types

import "time"

// Follow represents a following relationship between a leader and a follower.
type Follow struct {
	Leader   Username
	Follower Username
	Created  time.Time
}
