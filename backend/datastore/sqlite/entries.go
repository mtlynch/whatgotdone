package sqlite

import (
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

// ReadEntries returns all published entries matching the given filter.
func (d db) ReadEntries(filter datastore.EntryFilter) ([]types.JournalEntry, error) {
	whereClauses := []string{
		"is_draft=0",
	}
	values := []string{}
	if len(filter.ByUsers) != 0 {
		placeholders := strings.TrimSuffix(strings.Repeat("?,", len(filter.ByUsers)), ",")
		whereClauses = append(whereClauses, fmt.Sprintf("username IN (%s)", placeholders))
		for _, u := range filter.ByUsers {
			values = append(values, string(u))
		}
	}
	stmt, err := d.ctx.Prepare(fmt.Sprintf(`
		SELECT
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

	entries := []types.JournalEntry{}

	rows, err := stmt.Query(values)
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
