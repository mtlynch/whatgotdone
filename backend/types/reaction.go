package types

import "time"

// Reaction to an entity in What Got Done, such as a user liking a journal entry.
type Reaction struct {
	Username  Username  `json:"username"`
	Symbol    string    `json:"symbol"`
	Timestamp time.Time `json:"timestamp"`
}
