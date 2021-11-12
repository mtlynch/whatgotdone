package sqlite

import (
	"github.com/mtlynch/whatgotdone/backend/types"
)

// Users returns all the users who have published entries.
func (d db) Users() ([]types.Username, error) {
	return []types.Username{}, notImplementedError
}

// GetUserProfile returns profile information for the given user.
func (d db) GetUserProfile(username types.Username) (types.UserProfile, error) {
	/*stmt, err := d.ctx.Prepare("SELECT contents FROM entries WHERE id=?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var contents string
	err = stmt.QueryRow(id).Scan(&contents)
	if err != nil {
		return "", err
	}
	return contents, nil*/
	return types.UserProfile{}, notImplementedError
}

// SetUserProfile updates the given user's profile.
func (d db) SetUserProfile(username types.Username, profile types.UserProfile) error {
	/*_, err := d.ctx.Exec(`
	INSERT INTO entries(
		id,
		creation_time,
		contents)
	values(?,?,?)`, id, time.Now().Format(time.RFC3339), contents)
	return err*/
	return notImplementedError
}

// GetDraft returns an entry draft for the given user for the given date.
func (d db) GetDraft(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	return types.JournalEntry{}, notImplementedError
}

// InsertDraft saves an entry draft to the datastore, overwriting any existing
// draft with the same name and username.
func (d db) InsertDraft(username types.Username, j types.JournalEntry) error {
	return notImplementedError
}

// GetReactions retrieves reader reactions associated with a published entry.
func (d db) GetReactions(entryAuthor types.Username, entryDate types.EntryDate) ([]types.Reaction, error) {
	return []types.Reaction{}, notImplementedError
}

// AddReaction saves a reader reaction associated with a published entry,
// overwriting any existing reaction.
func (d db) AddReaction(entryAuthor types.Username, entryDate types.EntryDate, reaction types.Reaction) error {
	return notImplementedError
}

// InsertPageViews stores the count of pageviews for a given What Got Done route.
func (d db) InsertPageViews(path string, pageViews int) error {
	return notImplementedError
}

// GetPageViews retrieves the count of pageviews for a given What Got Done route.
func (d db) GetPageViews(path string) (int, error) {
	return 0, notImplementedError
}

// InsertFollow adds a following relationship to the datastore.
func (d db) InsertFollow(leader, follower types.Username) error {
	return notImplementedError
}

// DeleteFollow removes a following relationship from the datastore.
func (d db) DeleteFollow(leader, follower types.Username) error {
	return notImplementedError
}

// Followers returns all the users the specified user is following.
func (d db) Following(follower types.Username) ([]types.Username, error) {
	return []types.Username{}, notImplementedError
}
