package main

import (
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

//Cache struct
type Cache struct {
	client *redis.Client
	logger *log.Logger
}

//Settings contains the cache settings
type Settings struct {
	Address  string
	Password string
	DB       int
}

const cacheKeySignature = ":prev_command"

//Set sets  cache
func (cache *Cache) Set(key string, value interface{}, ttl time.Duration) {

	keyValue := key + cacheKeySignature
	err := cache.client.Set(keyValue, value, ttl).Err()
	if err != nil {
		cache.logger.Fatal("cache cannot be set for key :" + keyValue)
	}
}

//Get retrieves cache results
func (cache *Cache) Get(key string) string {
	result, err := cache.client.Get(key + cacheKeySignature).Result()
	if err != nil {
		cache.logger.Fatal("could not get cache for key :" + key)
	}
	return result
}

// NewCache creates a new cache instance
func NewCache(settings *Settings, logger *log.Logger) *Cache {

	redis := redis.NewClient(&redis.Options{
		Addr:     settings.Address,
		Password: settings.Password, // no password set
		DB:       settings.DB,       // use default DB
	})

	return &Cache{client: redis, logger: logger}

}
