package sqlite

import (
	"log"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetEntry returns the published entry for the given date.
func (d db) GetEntry(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	stmt, err := d.ctx.Prepare(`
		SELECT
			markdown, last_modified
		FROM
			journal_entries
		WHERE
			username=? AND
			date=? AND
			is_draft=0
		`)
	if err != nil {
		return types.JournalEntry{}, err
	}
	defer stmt.Close()

	var markdown string
	var lastModified string
	err = stmt.QueryRow(username, date).Scan(&markdown, &lastModified)
	if err != nil {
		return types.JournalEntry{}, err
	}

	t, err := parseDatetime(lastModified)
	if err != nil {
		return types.JournalEntry{}, err
	}

	return types.JournalEntry{
		Date:         date,
		LastModified: t.Format("2006-01-02T15:04:05Z"),
		Markdown:     markdown,
	}, nil
}

// GetEntries returns all published entries for the given user.
func (d db) GetEntries(username types.Username) ([]types.JournalEntry, error) {
	stmt, err := d.ctx.Prepare(`
		SELECT
			date,
			markdown,
			last_modified
		FROM
			journal_entries
		WHERE
			username=? AND
			is_draft=0
		`)
	if err != nil {
		return []types.JournalEntry{}, err
	}
	defer stmt.Close()

	entries := []types.JournalEntry{}

	rows, err := stmt.Query(username)
	for rows.Next() {
		var dateRaw string
		var markdown string
		var lastModified string
		err := rows.Scan(&dateRaw, &markdown, &lastModified)
		if err != nil {
			return []types.JournalEntry{}, err
		}

		date, err := parseDate(dateRaw)
		if err != nil {
			return []types.JournalEntry{}, err
		}

		t, err := parseDatetime(lastModified)
		if err != nil {
			return []types.JournalEntry{}, err
		}

		entries = append(entries, types.JournalEntry{
			Date:         types.EntryDate(date.Format("2006-01-02")),
			LastModified: t.Format("2006-01-02T15:04:05Z"),
			Markdown:     markdown,
		})
	}

	return entries, nil
}

// InsertEntry saves an entry to the datastore, overwriting any existing entry
// with the same name and username.
func (d db) InsertEntry(username types.Username, j types.JournalEntry) error {
	log.Printf("saving entry to datastore: %s -> %+v", username, j.Date)
	_, err := d.ctx.Exec(`
	INSERT INTO journal_entries(
		username,
		date,
		markdown,
		is_draft,
		last_modified)
	values(?,?,?,0,datetime('now'))`, username, j.Date, j.Markdown)
	return err
}
