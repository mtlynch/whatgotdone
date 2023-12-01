// Package datastore provides functionality for storing and retrieving
// persistent data. This package does not enforce access control, so it is the
// client's responsibility to enforce authentication and authorization before
// retrieving private data from the datastore.
package datastore

import (
	"fmt"
	"io"

	"github.com/mtlynch/whatgotdone/backend/types"
)

type EntryFilter struct {
	ByUsers   []types.Username
	MinLength int32
	Offset    int32
	Limit     int32
}

// Datastore represents the What Got Done datastore. It's responsible for
// storing and retrieving all persistent data (journal entries, journal drafts,
// reactions).
type Datastore interface {
	// GetAvatar returns avatar of the given user.
	GetAvatar(username types.Username) (io.Reader, error)
	InsertAvatar(username types.Username, avatar io.Reader, avatarWidth int) error
	DeleteAvatar(username types.Username) error
	// GetUserProfile returns profile information for the given user.
	GetUserProfile(username types.Username) (types.UserProfile, error)
	// SetUserProfile updates the given user's profile.
	SetUserProfile(username types.Username, profile types.UserProfile) error
	// GetEntry returns the published entry for the given date.
	GetEntry(username types.Username, date types.EntryDate) (types.JournalEntry, error)
	// ReadEntries returns all published entries matching the given filter.
	ReadEntries(filter EntryFilter) ([]types.JournalEntry, error)
	// GetDraft returns an entry draft for the given user for the given date.
	GetDraft(username types.Username, date types.EntryDate) (types.JournalEntry, error)
	// InsertEntry saves an entry to the datastore, overwriting any existing entry
	// with the same name and username.
	InsertEntry(username types.Username, j types.JournalEntry) error
	// DeleteEntry deletes an entry from the datastore.
	DeleteEntry(username types.Username, d types.EntryDate) error
	// InsertDraft saves an entry draft to the datastore, overwriting any existing
	// draft with the same name and username.
	InsertDraft(username types.Username, j types.JournalEntry) error
	// DeleteDraft deletes a draft from the datastore.
	DeleteDraft(username types.Username, date types.EntryDate) error
	// GetReactions retrieves reader reactions associated with a published entry.
	GetReactions(entryAuthor types.Username, entryDate types.EntryDate) ([]types.Reaction, error)
	// AddReaction saves a reader reaction associated with a published entry,
	// overwriting any existing reaction.
	AddReaction(entryAuthor types.Username, entryDate types.EntryDate, reaction types.Reaction) error
	// DeleteReaction removes a user's reaction to a published entry.
	DeleteReaction(entryAuthor types.Username, entryDate types.EntryDate, reactingUser types.Username) error
	// InsertFollow adds a following relationship to the datastore.
	InsertFollow(leader, follower types.Username) error
	// DeleteFollow removes a following relationship from the datastore.
	DeleteFollow(leader, follower types.Username) error
	// Followers returns all the users the specified user is following.
	Following(follower types.Username) ([]types.Username, error)
	// GetPreferences retrieves the user's preferences for using the site.
	GetPreferences(username types.Username) (types.Preferences, error)
	// SetPreferences saves the user's preferences for using the site.
	SetPreferences(username types.Username, prefs types.Preferences) error
}

// ErrAvatarNotFound occurs when no avatar exists for a user.
type ErrAvatarNotFound struct {
	Username types.Username
}

func (f ErrAvatarNotFound) Error() string {
	return fmt.Sprintf("no avatar found for username %s", f.Username)
}

// EntryNotFoundError occurs when no published exists for a user with a given date.
type EntryNotFoundError struct {
	Username types.Username
	Date     types.EntryDate
}

func (f EntryNotFoundError) Error() string {
	return fmt.Sprintf("Could not find published entry for user %s on date %s", f.Username, f.Date)
}

// DraftNotFoundError occurs when no draft exists for a user with a given date.
type DraftNotFoundError struct {
	Username types.Username
	Date     types.EntryDate
}

func (f DraftNotFoundError) Error() string {
	return fmt.Sprintf("Could not find draft entry for user %s on date %s", f.Username, f.Date)
}

// UserProfileNotFoundError occurs when no profile exists for the given
// username. The user might exist, but they have not submitted profile data.
type UserProfileNotFoundError struct {
	Username types.Username
}

func (f UserProfileNotFoundError) Error() string {
	return fmt.Sprintf("No user profile found for username %s", f.Username)
}

// PreferencesNotFoundError occurs when no profile exists for the given
// username. The user might exist, but they have not set preferences.
type PreferencesNotFoundError struct {
	Username types.Username
}

func (f PreferencesNotFoundError) Error() string {
	return fmt.Sprintf("No user preferences found for username %s", f.Username)
}
