package redis

import (
	redis "github.com/go-redis/redis/v7"

	"github.com/mtlynch/whatgotdone/backend/datastore"
)

// New creates a new Datastore instance implemented with Redis as the backend.
func New(datastoreAddress string) datastore.Datastore {
	addr := datastoreAddress
	if addr == "" {
		addr = "localhost:6379"
	}
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return client{
		redisClient: c,
	}
}
