package types

// Reaction to an entity in What Got Done, such as a user liking a journal entry.
type Reaction struct {
	Username  string `json:"username" firestore:"username,omitempty"`
	Symbol    string `json:"symbol" firestore:"symbol,omitempty"`
	Timestamp string `json:"timestamp" firestore:"timestamp,omitempty"`
}
