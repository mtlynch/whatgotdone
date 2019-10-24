package redis

import (
	redis "github.com/go-redis/redis/v7"

	"github.com/mtlynch/whatgotdone/backend/datastore"
)

// New creates a new Datastore instance implemented with Redis as the backend.
func New() datastore.Datastore {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client{
		redisClient: c,
	}
}
