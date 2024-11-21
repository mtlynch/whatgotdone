package sqlite

import (
	"database/sql"
	"log"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetPreferences retrieves the user's preferences for using the site.
func (d DB) GetPreferences(username types.Username) (types.Preferences, error) {
	var entryTemplate string
	err := d.ctx.QueryRow(`
		SELECT
				entry_template
		FROM
				user_preferences
		WHERE
				username = :username`,
		sql.Named("username", username)).Scan(&entryTemplate)

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
func (d DB) SetPreferences(username types.Username, prefs types.Preferences) error {
	log.Printf("saving preferences to datastore: %s -> %+v", username, prefs)
	_, err := d.ctx.Exec(`
		INSERT OR REPLACE INTO user_preferences(
				username,
				entry_template)
		values(:username, :entry_template)`,
		sql.Named("username", username),
		sql.Named("entry_template", prefs.EntryTemplate))
	return err
}
