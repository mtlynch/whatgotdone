package entries

import (
	"github.com/mtlynch/whatgotdone/backend/datastore"
)

type Reader interface {
	Recent(start, limit int) ([]RecentEntry, error)
	RecentFollowing(username string, start, limit int) ([]RecentEntry, error)
}

type defaultReader struct {
	datastore datastore.Datastore
}

func NewReader(ds datastore.Datastore) Reader {
	return defaultReader{
		datastore: ds,
	}
}
