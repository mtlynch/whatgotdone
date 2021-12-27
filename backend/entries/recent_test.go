package entries

import (
	"errors"
	"testing"

	"github.com/mtlynch/whatgotdone/backend/datastore/mock"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestRecentFailsWhenDatastoreFailsToRetrieveEntries(t *testing.T) {
	ms := mock.MockDatastore{
		Usernames: []types.Username{
			"bob",
		},
		ReadEntriesErr: errors.New("dummy error for MockDatastore.GetEntries()"),
	}
	r := defaultReader{
		store: &ms,
	}

	_, err := r.Recent(0, 20)
	if err == nil {
		t.Fatalf("Expected call to Recent to fail")
	}
}
