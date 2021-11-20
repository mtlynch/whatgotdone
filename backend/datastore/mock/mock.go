package mock

import (
	"errors"
	"sync"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// MockDatastore is a mock implementation of the datstore.Datastore interface
// for testing.
type MockDatastore struct {
	JournalEntries []types.JournalEntry
	JournalDrafts  []types.JournalEntry
	Usernames      []types.Username
	Reactions      []types.Reaction
	PageViewCounts []ga.PageViewCount
	UserProfile    types.UserProfile
	GetEntriesErr  error
	mu             sync.Mutex
}

func (ds *MockDatastore) Users() ([]types.Username, error) {
	return ds.Usernames, nil
}

func (ds *MockDatastore) GetUserProfile(username types.Username) (types.UserProfile, error) {
	return ds.UserProfile, nil
}

func (ds *MockDatastore) SetUserProfile(username types.Username, p types.UserProfile) error {
	ds.UserProfile = p
	return nil
}

func (ds *MockDatastore) GetEntry(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	if (username == types.Username("jimmy123")) && (date == types.EntryDate("2020-01-17")) {
		return types.JournalEntry{
			Markdown: "dummy journal content",
		}, nil
	}
	return types.JournalEntry{}, errors.New("mock not found")
}

func (ds *MockDatastore) GetEntries(username types.Username) ([]types.JournalEntry, error) {
	return ds.JournalEntries, ds.GetEntriesErr
}

func (ds *MockDatastore) GetDraft(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	if len(ds.JournalDrafts) > 0 {
		return ds.JournalDrafts[0], nil
	}
	return types.JournalEntry{}, datastore.DraftNotFoundError{
		Username: username,
		Date:     date,
	}
}

func (ds *MockDatastore) InsertEntry(username types.Username, j types.JournalEntry) error {
	return nil
}

func (ds *MockDatastore) InsertDraft(username types.Username, j types.JournalEntry) error {
	return nil
}

func (ds *MockDatastore) GetReactions(entryAuthor types.Username, entryDate types.EntryDate) ([]types.Reaction, error) {
	return ds.Reactions, nil
}

func (ds *MockDatastore) AddReaction(entryAuthor types.Username, entryDate types.EntryDate, reaction types.Reaction) error {
	ds.Reactions = append(ds.Reactions, reaction)
	return nil
}

func (ds *MockDatastore) DeleteReaction(entryAuthor types.Username, entryDate types.EntryDate, user types.Username) error {
	toDelete := -1
	for i, r := range ds.Reactions {
		if r.Username == user {
			toDelete = i
		}
	}
	if toDelete >= 0 {
		ds.Reactions = append(ds.Reactions[:toDelete], ds.Reactions[toDelete+1:]...)
	}
	return nil
}

func (ds *MockDatastore) InsertPageViews(path string, pageViews int) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.PageViewCounts = append(ds.PageViewCounts, ga.PageViewCount{
		Path:  path,
		Views: pageViews,
	})
	return nil
}

func (ds *MockDatastore) GetPageViews(path string) (int, error) {
	for _, pvc := range ds.PageViewCounts {
		if pvc.Path == path {
			return pvc.Views, nil
		}
	}
	return 0, errors.New("no pageview results found")
}

func (ds *MockDatastore) InsertFollow(leader, follower types.Username) error {
	return errors.New("MockDatastore does not implement InsertFollow")
}

func (ds *MockDatastore) DeleteFollow(leader, follower types.Username) error {
	return errors.New("MockDatastore does not implement DeleteFollow")
}

func (ds *MockDatastore) Following(follower types.Username) ([]types.Username, error) {
	return []types.Username{}, errors.New("MockDatastore does not implement Following")
}

func (ds *MockDatastore) GetPreferences(username types.Username) (types.Preferences, error) {
	return types.Preferences{}, datastore.PreferencesNotFoundError{
		Username: username,
	}
}

func (ds *MockDatastore) SetPreferences(username types.Username, prefs types.Preferences) error {
	return errors.New("MockDatastore does not implement SetPreferences")
}

func (ds *MockDatastore) Close() error {
	return nil
}
