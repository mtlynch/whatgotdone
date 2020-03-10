// Package datastore provides functionality for storing and retrieving
// persistent data. This package does not enforce access control, so it is the
// client's responsibility to enforce authentication and authorization before
// retrieving private data from the datastore.
package datastore

import (
	"fmt"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// Datastore represents the What Got Done datastore. It's responsible for
// storing and retrieving all persistent data (journal entries, journal drafts,
// reactions).
type Datastore interface {
	// Users returns all the users who have published entries.
	Users() ([]string, error)
	// GetUserProfile returns profile information for the given user.
	GetUserProfile(username string) (types.UserProfile, error)
	// SetUserProfile updates the given user's profile.
	SetUserProfile(username string, profile types.UserProfile) error
	// GetEntries returns all published entries for the given user.
	GetEntries(username string) ([]types.JournalEntry, error)
	// GetDraft returns an entry draft for the given user for the given date.
	GetDraft(username string, date string) (types.JournalEntry, error)
	// InsertEntry saves an entry to the datastore, overwriting any existing entry
	// with the same name and username.
	InsertEntry(username string, j types.JournalEntry) error
	// InsertDraft saves an entry draft to the datastore, overwriting any existing
	// draft with the same name and username.
	InsertDraft(username string, j types.JournalEntry) error
	// GetReactions retrieves reader reactions associated with a published entry.
	GetReactions(entryAuthor string, entryDate string) ([]types.Reaction, error)
	// AddReaction saves a reader reaction associated with a published entry,
	// overwriting any existing reaction.
	AddReaction(entryAuthor string, entryDate string, reaction types.Reaction) error
	// InsertPageViews stores the count of pageviews for a given What Got Done route.
	InsertPageViews(path string, pageViews int) error
	// GetPageViews retrieves the count of pageviews for a given What Got Done route.
	GetPageViews(path string) (int, error)
	// InsertFollow adds a following relationship to the datastore.
	InsertFollow(leader, follower string) error
	// DeleteFollow removes a following relationship from the datastore.
	DeleteFollow(leader, follower string) error
	// Followers returns all the users the specified user is following.
	Following(follower string) ([]string, error)
	// GetPreferences retrieves the user's preferences for using the site.
	GetPreferences(username string) (types.Preferences, error)
	// SetPreferences saves the user's preferences for using the site.
	SetPreferences(username string, prefs types.Preferences) error
}

// DraftNotFoundError occurs when no draft exists for a user with a given date.
type DraftNotFoundError struct {
	Username string
	Date     string
}

func (f DraftNotFoundError) Error() string {
	return fmt.Sprintf("Could not find draft entry for user %s on date %s", f.Username, f.Date)
}

// UserProfileNotFoundError occurs when no profile exists for the given
// username. The user might exist, but they have not submitted profile data.
type UserProfileNotFoundError struct {
	Username string
}

func (f UserProfileNotFoundError) Error() string {
	return fmt.Sprintf("No user profile found for username %s", f.Username)
}

// PageViewsNotFoundError occurs when no page view data is present in the
// datastore for the given URL path.
type PageViewsNotFoundError struct {
	Path string
}

func (f PageViewsNotFoundError) Error() string {
	return fmt.Sprintf("No page view count found for path %s", f.Path)
}

type FollowAlreadyExistsError struct {
	Leader   string
	Follower string
}

func (f FollowAlreadyExistsError) Error() string {
	return fmt.Sprintf("Cannot create a follow from %s -> %s because the follow already exists", f.Follower, f.Leader)
}

// PreferencesNotFoundError occurs when no profile exists for the given
// username. The user might exist, but they have not set preferences.
type PreferencesNotFoundError struct {
	Username string
}

func (f PreferencesNotFoundError) Error() string {
	return fmt.Sprintf("No user preferences found for username %s", f.Username)
}
