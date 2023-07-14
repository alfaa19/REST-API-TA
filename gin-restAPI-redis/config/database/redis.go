package config

import "github.com/go-redis/redis/v8"

var RedisDB *redis.Client

func ConnectRedis() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

}
