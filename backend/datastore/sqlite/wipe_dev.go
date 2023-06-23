//go:build dev || staging

package sqlite

import "log"

// Clear drops all SQLite tables. Only for testing.
func (db DB) Clear() {
	log.Printf("clearing all SQLite tables")
	tables := []string{
		"user_preferences",
		"user_profiles",
		"journal_entries",
		"follows",
		"entry_reactions",
	}
	for _, tbl := range tables {
		if _, err := db.ctx.Exec("DELETE FROM " + tbl); err != nil {
			log.Fatalf("failed to delete table %s: %v", tbl, err)
		}
	}
}
