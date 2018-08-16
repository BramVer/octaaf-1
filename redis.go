package main

import (
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/gobuffalo/envy"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
)

// Redis is is the redis client struct
var Redis *redis.Client

// Codec is the fast redis cacher
var Cache *cache.Codec

func initRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     envy.Get("REDIS_URI", "localhost:6379"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	log.Info("Established Redis connection")

	Cache = &cache.Codec{
		Redis: Redis,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	log.Info("Established Redis cache")
}
