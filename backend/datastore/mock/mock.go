package mock

import (
	"errors"
	"io"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// MockDatastore is a mock implementation of the datstore.Datastore interface
// for testing.
type MockDatastore struct {
	JournalEntries    []types.JournalEntry
	JournalDrafts     []types.JournalEntry
	Usernames         []types.Username
	Reactions         map[types.Username]map[types.EntryDate][]types.Reaction
	UserFollows       map[types.Username][]types.Username
	UserPreferences   map[types.Username]types.Preferences
	ForwardingAddress map[types.Username]types.ForwardingAddress
	UserProfile       types.UserProfile
	ReadEntriesErr    error
}

func (ds MockDatastore) GetAvatar(username types.Username) (io.Reader, error) {
	return nil, errors.New("not implemented")
}

func (ds MockDatastore) InsertAvatar(username types.Username, avatarReader io.Reader, avatarWidth int) error {
	return errors.New("not implemented")
}

func (ds MockDatastore) DeleteAvatar(username types.Username) error {
	return errors.New("not implemented")
}

func (ds *MockDatastore) GetUserProfile(username types.Username) (types.UserProfile, error) {
	return ds.UserProfile, nil
}

func (ds *MockDatastore) SetUserProfile(username types.Username, p types.UserProfile) error {
	ds.UserProfile = p
	return nil
}

func (ds *MockDatastore) GetEntry(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	for _, entry := range ds.JournalEntries {
		if entry.Author == username && entry.Date == date {
			return entry, nil
		}
	}
	return types.JournalEntry{}, errors.New("mock journal entry not found")
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

func (ds *MockDatastore) DeleteDraft(username types.Username, date types.EntryDate) error {
	for i, draft := range ds.JournalDrafts {
		if draft.Author == username && draft.Date == date {
			ds.JournalDrafts = append(ds.JournalDrafts[:i], ds.JournalDrafts[i+1:]...)
			return nil
		}
	}

	return nil
}

func (ds *MockDatastore) InsertEntry(username types.Username, j types.JournalEntry) error {
	return nil
}

func (ds *MockDatastore) DeleteEntry(username types.Username, date types.EntryDate) error {
	for i, entry := range ds.JournalEntries {
		if entry.Author == username && entry.Date == date {
			ds.JournalEntries = append(ds.JournalEntries[:i], ds.JournalEntries[i+1:]...)
			return nil
		}
	}

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

func (ds *MockDatastore) GetForwardingAddress(username types.Username) (types.ForwardingAddress, error) {
	address, ok := ds.ForwardingAddress[username]
	if !ok {
		return types.ForwardingAddress(""), datastore.ForwardingAddressNotFoundError{
			Username: username,
		}
	}
	return address, nil
}

func (ds *MockDatastore) SetForwardingAddress(username types.Username, address types.ForwardingAddress) error {
	if ds.ForwardingAddress == nil {
		ds.ForwardingAddress = make(map[types.Username]types.ForwardingAddress)
	}
	ds.ForwardingAddress[username] = address
	return nil
}

func (ds *MockDatastore) DeleteForwardingAddress(username types.Username) error {
	if ds.ForwardingAddress != nil {
		delete(ds.ForwardingAddress, username)
	}
	return nil
}

func (ds *MockDatastore) Close() error {
	return nil
}
