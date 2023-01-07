package sqlite

import (
	"database/sql"
	"log"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetDraft returns an entry draft for the given user for the given date.
func (d DB) GetDraft(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	stmt, err := d.ctx.Prepare(`
			SELECT
				markdown,
				last_modified
			FROM
				journal_entries
			WHERE
				username=? AND
				date=? AND
				is_draft=1
			`)
	if err != nil {
		return types.JournalEntry{}, err
	}
	defer stmt.Close()

	var markdown string
	var lastModified string
	err = stmt.QueryRow(username, date).Scan(&markdown, &lastModified)
	if err == sql.ErrNoRows {
		return types.JournalEntry{}, datastore.DraftNotFoundError{
			Username: username,
			Date:     date,
		}
	} else if err != nil {
		return types.JournalEntry{}, err
	}

	t, err := parseDatetime(lastModified)
	if err != nil {
		return types.JournalEntry{}, err
	}

	return types.JournalEntry{
		Author:       username,
		Date:         date,
		LastModified: t,
		Markdown:     types.EntryContent(markdown),
	}, nil
}

// InsertDraft saves an entry draft to the datastore, overwriting any existing
// entry with the same name and username.
func (d DB) InsertDraft(username types.Username, j types.JournalEntry) error {
	log.Printf("saving draft to datastore: %s -> %+v (%d characters)", username, j.Date, len(j.Markdown))
	_, err := d.ctx.Exec(`
	INSERT OR REPLACE INTO journal_entries(
		username,
		date,
		markdown,
		is_draft,
		last_modified)
	values(?,?,?,1,strftime('%Y-%m-%d %H:%M:%SZ', 'now', 'utc'))`, username, j.Date, j.Markdown)
	return err
}

// DeleteDraft deletes a draft from the datastore.
func (d DB) DeleteDraft(username types.Username, date types.EntryDate) error {
	log.Printf("deleting draft from datastore: %s -> %+v", username, date)
	_, err := d.ctx.Exec(`
	DELETE FROM
		journal_entries
	WHERE
		username=? AND
		date=? AND
		is_draft=1
	`, username, date)
	return err
}
