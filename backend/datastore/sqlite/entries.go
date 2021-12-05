package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetEntry returns the published entry for the given date.
func (d db) GetEntry(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	stmt, err := d.ctx.Prepare(`
		SELECT
			markdown,
			last_modified
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
	if err == sql.ErrNoRows {
		return types.JournalEntry{}, datastore.EntryNotFoundError{
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

// ReadEntries returns all published entries matching the given filter.
func (d db) ReadEntries(filter datastore.EntryFilter) ([]types.JournalEntry, error) {
	whereClauses := []string{
		"is_draft=0",
	}
	var values []interface{}
	if len(filter.ByUsers) != 0 {
		placeholders := strings.TrimSuffix(strings.Repeat("?,", len(filter.ByUsers)), ",")
		whereClauses = append(whereClauses, fmt.Sprintf("username IN (%s)", placeholders))
		for _, u := range filter.ByUsers {
			values = append(values, string(u))
		}
	}

	if filter.MinLength != 0 {
		whereClauses = append(whereClauses, "LENGTH(markdown) > ?")
		values = append(values, filter.MinLength)
	}

	stmt, err := d.ctx.Prepare(fmt.Sprintf(`
		SELECT
			username,
			date,
			markdown,
			last_modified
		FROM
			journal_entries
		WHERE
		  %s
		`, strings.Join(whereClauses, " AND ")))
	if err != nil {
		return []types.JournalEntry{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(values...)
	if err != nil {
		return []types.JournalEntry{}, err
	}

	entries := []types.JournalEntry{}
	for rows.Next() {
		var usernameRaw string
		var dateRaw string
		var markdown string
		var lastModified string
		err := rows.Scan(&usernameRaw, &dateRaw, &markdown, &lastModified)
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
			Author:       types.Username(usernameRaw),
			Date:         types.EntryDate(date.Format("2006-01-02")),
			LastModified: t,
			Markdown:     types.EntryContent(markdown),
		})
	}

	return entries, nil
}

// InsertEntry saves an entry to the datastore, overwriting any existing entry
// with the same name and username.
func (d db) InsertEntry(username types.Username, j types.JournalEntry) error {
	log.Printf("saving entry to datastore: %s -> %+v (%d characters)", username, j.Date, len(j.Markdown))
	_, err := d.ctx.Exec(`
	INSERT OR REPLACE INTO journal_entries(
		username,
		date,
		markdown,
		is_draft,
		last_modified)
	values(?,?,?,0,strftime('%Y-%m-%d %H:%M:%SZ', 'now', 'utc'))`, username, j.Date, j.Markdown)
	return err
}

// DeleteDraft deletes a draft from the datastore.
func (d db) DeleteEntry(username types.Username, date types.EntryDate) error {
	log.Printf("deleting entry from datastore: %s -> %+v", username, date)
	_, err := d.ctx.Exec(`
	DELETE FROM
		journal_entries
	WHERE
		username=? AND
		date=? AND
		is_draft=0
	`, username, date)
	return err
}
