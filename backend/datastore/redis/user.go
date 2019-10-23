package redis

import (
	"encoding/json"
	"fmt"
	"strings"

	redis "github.com/go-redis/redis/v7"
	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// Users returns all the users who have published entries.
func (c client) Users() ([]string, error) {
	keyPattern := formatEntryKey("*", "*")
	keys, err := c.redisClient.Keys(keyPattern).Result()
	if err != nil {
		return []string{}, err
	}
	users := map[string]bool{}
	for _, k := range keys {
		username := usernameFromEntryKey(k)
		users[username] = true
	}

	uniqueUsers := []string{}
	for u := range users {
		uniqueUsers = append(uniqueUsers, u)
	}

	return uniqueUsers, nil
}

// UserProfile returns profile information about the given user.
func (c client) GetUserProfile(username string) (types.UserProfile, error) {
	key := formatUserProfileKey(username)
	profileJSON, err := c.redisClient.Get(key).Result()
	if err == redis.Nil {
		return types.UserProfile{}, datastore.UserProfileNotFoundError{
			Username: username,
		}
	} else if err != nil {
		return types.UserProfile{}, err
	}

	var p types.UserProfile
	if err = json.Unmarshal([]byte(profileJSON), &p); err != nil {
		return types.UserProfile{}, err
	}

	return p, nil
}

// SetUserProfile updates the given user's profile or creates a new profile for
// the user.
func (c client) SetUserProfile(username string, p types.UserProfile) error {
	profileJSONBytes, err := json.Marshal(p)
	if err != nil {
		return err
	}
	key := formatUserProfileKey(username)
	profileJSON := string(profileJSONBytes)
	err = c.redisClient.Set(key, profileJSON, 0).Err()
	if err != nil {
		return err
	}
	return err
}

func usernameFromEntryKey(entryKey string) string {
	p := strings.SplitN(entryKey, ":", 3)
	return p[1]
}

func formatUserProfileKey(username string) string {
	return fmt.Sprintf("userProfiles:%s", username)
}
