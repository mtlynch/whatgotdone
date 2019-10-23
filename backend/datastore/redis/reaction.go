package redis

import (
	"encoding/json"
	"fmt"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetReactions retrieves reader reactions associated with a published entry.
func (c client) GetReactions(entryAuthor string, entryDate string) ([]types.Reaction, error) {
	keyPattern := formatReactionKey(entryAuthor, entryDate, "*")
	keys, err := c.redisClient.Keys(keyPattern).Result()
	if err != nil {
		return []types.Reaction{}, err
	}
	if len(keys) == 0 {
		return []types.Reaction{}, nil
	}

	reactionsRaw, err := c.redisClient.MGet(keys...).Result()
	if err != nil {
		return []types.Reaction{}, err
	}
	reactions := []types.Reaction{}
	for _, reactionJSON := range reactionsRaw {
		var r types.Reaction
		if err = json.Unmarshal([]byte(reactionJSON.(string)), &r); err != nil {
			return []types.Reaction{}, err
		}
		reactions = append(reactions, r)
	}
	return reactions, nil
}

// AddReaction saves a reader reaction associated with a published entry,
// overwriting any existing reaction.
func (c client) AddReaction(entryAuthor string, entryDate string, reaction types.Reaction) error {
	reactionJSONBytes, err := json.Marshal(reaction)
	if err != nil {
		return err
	}
	reactionJSON := string(reactionJSONBytes)
	key := formatReactionKey(entryAuthor, entryDate, reaction.Username)
	err = c.redisClient.Set(key, reactionJSON, 0).Err()
	if err != nil {
		return err
	}
	return err
}

func formatReactionKey(entryAuthor string, entryDate string, reactingUser string) string {
	return fmt.Sprintf("reactions:%s:%s:%s", entryAuthor, entryDate, reactingUser)
}
