package types

import "time"

// Reaction to an entity in What Got Done, such as a user liking a journal entry.
type Reaction struct {
	Username string `json:"username" firestore:"username,omitempty"`
	Symbol   string `json:"symbol" firestore:"symbol,omitempty"`
	// Timestamp is deprecated. It is pending deletion once the database is
	// updated to populate CreationTime in all the stored values.
	Timestamp    string    `json:"timestamp" firestore:"timestamp,omitempty"`
	CreationTime time.Time `json:"creationTime" firestore:"timestamp"`
}
