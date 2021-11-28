package mock

import (
	"errors"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// MockDatastore is a mock implementation of the datstore.Datastore interface
// for testing.
type MockDatastore struct {
	JournalEntries  []types.JournalEntry
	JournalDrafts   []types.JournalEntry
	Usernames       []types.Username
	Reactions       map[types.Username]map[types.EntryDate][]types.Reaction
	UserFollows     map[types.Username][]types.Username
	UserPreferences map[types.Username]types.Preferences
	PageViewCounts  []ga.PageViewCount
	UserProfile     types.UserProfile
	ReadEntriesErr  error
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

func (ds *MockDatastore) ReadEntries(filter datastore.EntryFilter) ([]types.JournalEntry, error) {
	if ds.ReadEntriesErr != nil {
		return []types.JournalEntry{}, ds.ReadEntriesErr
	}
	var entries []types.JournalEntry
	for _, e := range ds.JournalEntries {
		if len(filter.ByUsers) == 0 {
			entries = append(entries, e)
			continue
		}
		for _, u := range filter.ByUsers {
			if e.Author == u {
				entries = append(entries, e)
			}
		}
	}
	return entries, nil
}

func (ds *MockDatastore) GetDraft(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	for _, d := range ds.JournalDrafts {
		if d.Date == date {
			return d, nil
		}
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
	for username, reactionsByDate := range ds.Reactions {
		if entryAuthor != username {
			continue
		}
		for date, reactions := range reactionsByDate {
			if date != entryDate {
				continue
			}
			return reactions, nil
		}
	}
	return []types.Reaction{}, nil
}

func (ds *MockDatastore) AddReaction(entryAuthor types.Username, entryDate types.EntryDate, reaction types.Reaction) error {
	reactionsToAuthor, ok := ds.Reactions[entryAuthor]
	if !ok {
		reactionsToAuthor = map[types.EntryDate][]types.Reaction{}
		ds.Reactions[entryAuthor] = reactionsToAuthor
	}
	reactionsToAuthor[entryDate] = append(reactionsToAuthor[entryDate], reaction)
	return nil
}

func (ds *MockDatastore) DeleteReaction(entryAuthor types.Username, entryDate types.EntryDate, user types.Username) error {
	for username, reactionsByDate := range ds.Reactions {
		if entryAuthor != username {
			continue
		}
		toDelete := -1
		for i, r := range reactionsByDate[entryDate] {
			if r.Username == user {
				toDelete = i
				break
			}
		}
		if toDelete >= 0 {
			ds.Reactions[username][entryDate] = append(ds.Reactions[username][entryDate][:toDelete], ds.Reactions[username][entryDate][toDelete+1:]...)
			break
		}
	}
	return nil
}

func (ds *MockDatastore) InsertPageViews(pvc []ga.PageViewCount) error {
	ds.PageViewCounts = pvc
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
	return ds.UserFollows[follower], nil
}

func (ds *MockDatastore) GetPreferences(username types.Username) (types.Preferences, error) {
	prefs, ok := ds.UserPreferences[username]
	if !ok {
		return types.Preferences{}, datastore.PreferencesNotFoundError{
			Username: username,
		}
	}
	return prefs, nil
}

func (ds *MockDatastore) SetPreferences(username types.Username, prefs types.Preferences) error {
	return errors.New("MockDatastore does not implement SetPreferences")
}

func (ds *MockDatastore) Close() error {
	return nil
}
