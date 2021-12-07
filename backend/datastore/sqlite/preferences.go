package sqlite

import (
	"database/sql"
	"log"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetPreferences retrieves the user's preferences for using the site.
func (d db) GetPreferences(username types.Username) (types.Preferences, error) {
	stmt, err := d.ctx.Prepare(`
	SELECT
		entry_template
	FROM
		user_preferences
	WHERE
		username=?`)
	if err != nil {
		return types.Preferences{}, err
	}
	defer stmt.Close()

	var entryTemplate string
	err = stmt.QueryRow(username).Scan(&entryTemplate)
	if err == sql.ErrNoRows {
		return types.Preferences{}, datastore.PreferencesNotFoundError{
			Username: username,
		}
	} else if err != nil {
		return types.Preferences{}, err
	}

	return types.Preferences{
		EntryTemplate: types.EntryContent(entryTemplate),
	}, nil
}

// SetPreferences saves the user's preferences for using the site.
func (d db) SetPreferences(username types.Username, prefs types.Preferences) error {
	log.Printf("saving preferences to datastore: %s -> %+v", username, prefs)
	_, err := d.ctx.Exec(`
	INSERT OR REPLACE INTO user_preferences(
		username,
		entry_template)
	values(?,?)`, username, prefs.EntryTemplate)
	return err
}
