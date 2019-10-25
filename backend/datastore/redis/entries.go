package redis

import (
	"encoding/json"
	"fmt"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetEntries returns all published entries for the given user.
func (c client) GetEntries(username string) ([]types.JournalEntry, error) {
	keyPattern := formatEntryKey(username, "*")
	keys, err := c.redisClient.Keys(keyPattern).Result()
	if err != nil {
		return []types.JournalEntry{}, err
	}
	if len(keys) == 0 {
		return []types.JournalEntry{}, nil
	}

	entriesJSON, err := c.redisClient.MGet(keys...).Result()
	if err != nil {
		return []types.JournalEntry{}, err
	}
	entries := []types.JournalEntry{}
	for _, entryJSON := range entriesJSON {
		var j types.JournalEntry
		if err = json.Unmarshal([]byte(entryJSON.(string)), &j); err != nil {
			return []types.JournalEntry{}, err
		}
		entries = append(entries, j)
	}
	return entries, nil
}

// InsertEntry saves an entry to the datastore, overwriting any existing entry
// with the same name and username.
func (c client) InsertEntry(username string, j types.JournalEntry) error {
	entryJSONBytes, err := json.Marshal(j)
	if err != nil {
		return err
	}
	entryJSON := string(entryJSONBytes)
	key := formatEntryKey(username, j.Date)
	err = c.redisClient.Set(key, entryJSON, 0).Err()
	if err != nil {
		return err
	}
	return err
}

func formatEntryKey(username string, date string) string {
	return fmt.Sprintf("entries:%s:%s", username, date)
}
