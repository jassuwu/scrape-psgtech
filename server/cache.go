package server

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redisClient *redis.Client
	bgCtx       context.Context
}

func NewCache(Addr string) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: "", // Setting no password
		DB:       0,  // Use default DB
	})
	ctx := context.Background()

	if err := client.FlushDB(ctx).Err(); err != nil {
		log.Fatalln("Failed to flush Redis DB: ", err)
	}

	log.Println("ALL CACHE KEYS INVALIDATED")
	return &Cache{
		redisClient: client,
		bgCtx:       ctx,
	}
}

func (c *Cache) Set(key string, val any) (bool, error) {
	data, err := json.Marshal(val)
	if err != nil {
		log.Println("Error marshalling data for caching: ", err)
		return false, err
	}
	err = c.redisClient.Set(c.bgCtx, key, data, 0).Err()
	if err != nil {
		log.Println("Error setting cache: ", err)
		return false, err
	}
	log.Println("Cached results for search: ", key)
	return true, nil
}

func (c *Cache) Get(key string) ([]RankedDocument, error) {
	data, err := c.redisClient.Get(c.bgCtx, key).Result()
	// redis.Nil indicates that the key doesn't exist AKA cache miss
	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		log.Println("Error retrieving from cache: ", err)
		return nil, err
	}

	results := []RankedDocument{}
	err = json.Unmarshal([]byte(data), &results)
	if err != nil {
		log.Println("Error unmarshalling cached data: ", err)
		return nil, err
	}
	return results, nil
}
