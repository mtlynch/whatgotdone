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
func (d DB) GetEntry(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	var markdown string
	var lastModified string
	err := d.ctx.QueryRow(`
		SELECT
			markdown,
			last_modified
		FROM
			journal_entries
		WHERE
			username=:username AND
			date=:date AND
			is_draft=0
		`, sql.Named("username", username), sql.Named("date", date)).Scan(&markdown, &lastModified)

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
func (d DB) ReadEntries(filter datastore.EntryFilter) ([]types.JournalEntry, error) {
	whereClauses := []string{
		"is_draft=0",
	}
	var namedArgs []sql.NamedArg
	if len(filter.ByUsers) != 0 {
		// Construct the placeholders dynamically based on the number of users.
		placeholders := []string{}
		for i := range filter.ByUsers {
			placeholders = append(placeholders, fmt.Sprintf(":user%d", i))
			namedArgs = append(namedArgs, sql.Named(fmt.Sprintf("user%d", i), filter.ByUsers[i]))
		}

		whereClauses = append(whereClauses, fmt.Sprintf("username IN (%s)", strings.Join(placeholders, ",")))

	}

	if filter.MinLength != 0 {
		whereClauses = append(whereClauses, "LENGTH(markdown) > :minLength")
		namedArgs = append(namedArgs, sql.Named("minLength", filter.MinLength))
	}

	limitClause := ""
	if filter.Limit != 0 {
		limitClause = "LIMIT :limit"
		namedArgs = append(namedArgs, sql.Named("limit", filter.Limit))

	}

	offsetClause := ""
	if filter.Offset != 0 {
		offsetClause = "OFFSET :offset"
		namedArgs = append(namedArgs, sql.Named("offset", filter.Offset))
	}

	query := fmt.Sprintf(`
		SELECT
			username,
			date,
			markdown,
			last_modified
		FROM
			journal_entries
		WHERE
			%s
		ORDER BY
			date DESC,
			last_modified DESC
			%s
			%s
		`, strings.Join(whereClauses, " AND "), limitClause, offsetClause)

	args := make([]any, len(namedArgs))
	for i, arg := range namedArgs {
		args[i] = arg
	}

	rows, err := d.ctx.Query(query, args...)
	if err != nil {
		return []types.JournalEntry{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close SQLite rows: %v", err)
		}
	}()

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
func (d DB) InsertEntry(username types.Username, j types.JournalEntry) error {
	log.Printf("saving entry to datastore: %s -> %+v (%d characters)", username, j.Date, len(j.Markdown))
	_, err := d.ctx.Exec(`
		INSERT OR REPLACE INTO journal_entries(
			username,
			date,
			markdown,
			is_draft,
			last_modified)
		values(:username,:date,:markdown,0,strftime('%Y-%m-%d %H:%M:%SZ', 'now', 'utc'))`,
		sql.Named("username", username),
		sql.Named("date", j.Date),
		sql.Named("markdown", j.Markdown))
	return err
}

// DeleteEntry deletes an entry from the datastore.
func (d DB) DeleteEntry(username types.Username, date types.EntryDate) error {
	log.Printf("deleting entry from datastore: %s -> %+v", username, date)
	_, err := d.ctx.Exec(`
		DELETE FROM
			journal_entries
		WHERE
			username=:username AND
			date=:date AND
			is_draft=0
		`, sql.Named("username", username), sql.Named("date", date))
	return err
}
