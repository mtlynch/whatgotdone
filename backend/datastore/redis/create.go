package redis

import (
	redis "github.com/go-redis/redis/v7"

	"github.com/mtlynch/whatgotdone/backend/datastore"
)

// New creates a new Datastore instance.
func New() datastore.Datastore {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client{
		redisClient: c,
	}
}
