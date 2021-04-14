package entries

import (
	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type Reader interface {
	Recent() ([]types.JournalEntry, error)
}

type defaultReader struct {
	datastore datastore.Datastore
}

func NewReader(ds datastore.Datastore) Reader {
	return defaultReader{
		datastore: ds,
	}
}
