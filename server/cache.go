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
	return &Cache{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     Addr,
			Password: "", // Setting no password
			DB:       0,  // Use default DB
		}),
		bgCtx: context.Background(),
	}
}

func (c *Cache) Set(key string, val any) (bool, error) {
	err := c.redisClient.HSet(c.bgCtx, key, val).Err()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func (c *Cache) Get(key string) (map[string]string, error) {
	results, err := c.redisClient.HGetAll(c.bgCtx, key).Result()
	if err != nil {
		return nil, err
	}
	return results, nil
}
