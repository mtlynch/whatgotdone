package types

// Reaction to an entity in What Got Done, such as a user liking a journal entry.
type Reaction struct {
	Username  string `json:"username" firestore:"username,omitempty"`
	Reaction  string `json:"reaction" firestore:"reaction,omitempty"`
	Timestamp string `json:"timestamp" firestore:"timestamp,omitempty"`
}
