package handlers

import (
	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/datastore/sqlite"
)

// New creates a new sqlite-based Datastore instance.
func newDatastore() datastore.Datastore {
	return sqlite.New()
}
