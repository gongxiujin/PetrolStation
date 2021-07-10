package application

import (
	redis "github.com/go-redis/redis/v8"
	"time"
)

var Redis *redis.Client


func InitRedisClient() *redis.Client{
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	return rdb
}