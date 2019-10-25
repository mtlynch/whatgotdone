package redis

import (
	"encoding/json"
	"fmt"

	redis "github.com/go-redis/redis/v7"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetDraft returns an entry draft for the given user for the given date.
func (c client) GetDraft(username string, date string) (types.JournalEntry, error) {
	key := formatDraftKey(username, date)
	draftJSON, err := c.redisClient.Get(key).Result()
	if err == redis.Nil {
		return types.JournalEntry{}, datastore.DraftNotFoundError{
			Username: username,
			Date:     date,
		}
	} else if err != nil {
		return types.JournalEntry{}, err
	}

	var j types.JournalEntry
	if err = json.Unmarshal([]byte(draftJSON), &j); err != nil {
		return types.JournalEntry{}, err
	}

	return j, nil
}

// InsertDraft saves an entry draft to the datastore, overwriting any existing
// entry with the same name and username.
func (c client) InsertDraft(username string, j types.JournalEntry) error {
	key := formatDraftKey(username, j.Date)

	draftJSONBytes, err := json.Marshal(j)
	if err != nil {
		return err
	}
	draftJSON := string(draftJSONBytes)
	err = c.redisClient.Set(key, draftJSON, 0).Err()
	if err != nil {
		return err
	}
	return err
}

func formatDraftKey(username string, date string) string {
	return fmt.Sprintf("journalDrafts:%s:%s", username, date)
}
