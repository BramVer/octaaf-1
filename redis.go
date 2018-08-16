package main

import (
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/gobuffalo/envy"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
)

func getRedis() *redis.Client {
	defer log.Info("Established Redis connection")
	return redis.NewClient(&redis.Options{
		Addr:     envy.Get("REDIS_URI", "localhost:6379"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func getCodec() *cache.Codec {
	defer log.Info("Established Redis cache")
	return &cache.Codec{
		Redis: state.Redis,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}
