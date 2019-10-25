// +build !redis

package handlers

import (
	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/datastore/firestore"
)

// New creates a new firestore-based Datastore instance.
func newDatastore(datastoreAddr string) datastore.Datastore {
	return firestore.New(datastoreAddr)
}
