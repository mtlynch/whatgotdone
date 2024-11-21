package sqlite

import (
	"database/sql"
	"log"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetReactions retrieves reader reactions associated with a published entry.
func (d DB) GetReactions(entryAuthor types.Username, entryDate types.EntryDate) ([]types.Reaction, error) {
	rows, err := d.ctx.Query(`
		SELECT
				reacting_user,
				reaction,
				timestamp
		FROM
				entry_reactions
		WHERE
				entry_author = :entry_author AND
				entry_date = :entry_date`,
		sql.Named("entry_author", entryAuthor),
		sql.Named("entry_date", entryDate))
	if err != nil {
		return []types.Reaction{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close SQLite rows: %v", err)
		}
	}()

	reactions := []types.Reaction{}
	for rows.Next() {
		var user string
		var reaction string
		var timestampRaw string
		err := rows.Scan(&user, &reaction, &timestampRaw)
		if err != nil {
			return []types.Reaction{}, err
		}

		t, err := parseDatetime(timestampRaw)
		if err != nil {
			return []types.Reaction{}, err
		}

		reactions = append(reactions, types.Reaction{
			Username:  types.Username(user),
			Symbol:    reaction,
			Timestamp: t,
		})
	}

	return reactions, nil
}

// AddReaction saves a reader reaction associated with a published entry,
// overwriting any existing reaction.
func (d DB) AddReaction(entryAuthor types.Username, entryDate types.EntryDate, reaction types.Reaction) error {
	log.Printf("saving reaction to datastore: %s to %s/%s: [%s]", reaction.Username, entryAuthor, entryDate, reaction.Symbol)
	_, err := d.ctx.Exec(`
		INSERT OR REPLACE INTO entry_reactions(
				entry_author,
				entry_date,
				reacting_user,
				reaction,
				timestamp)
		values(:entry_author, :entry_date, :reacting_user, :reaction, strftime('%Y-%m-%d %H:%M:%SZ', 'now', 'utc'))`,
		sql.Named("entry_author", entryAuthor),
		sql.Named("entry_date", entryDate),
		sql.Named("reacting_user", reaction.Username),
		sql.Named("reaction", reaction.Symbol))
	return err
}

// DeleteReaction removes a user's reaction to a published entry.
func (d DB) DeleteReaction(entryAuthor types.Username, entryDate types.EntryDate, reactingUser types.Username) error {
	log.Printf("deleting reaction from datastore: %s to %s/%s", reactingUser, entryAuthor, entryDate)
	_, err := d.ctx.Exec(`
		DELETE FROM
				entry_reactions
		WHERE
				entry_author = :entry_author AND
				entry_date = :entry_date AND
				reacting_user = :reacting_user`,
		sql.Named("entry_author", entryAuthor),
		sql.Named("entry_date", entryDate),
		sql.Named("reacting_user", reactingUser))
	return err
}
