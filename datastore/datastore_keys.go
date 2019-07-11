package datastore

const (
	entriesRootKey      = "journalEntries"
	perUserEntriesKey   = "entries"
	draftsRootKey       = "journalDrafts"
	perUserDraftsKey    = "drafts"
	reactionsRootKey    = "reactions"
	secretsRootKey      = "secrets"
	secretUserKitDocKey = "userKitKey"
)

var (
	// perUserReactionsKey needs to be a variable because in our e2e tests, we
	// reassign the variable at compile time to ensure empty reactions in the
	// datastore.
	perUserReactionsKey = "perUserReactions"
)
