package bible

import (
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache interface {
	//Get gets a string from cache using provided key.
	Get(key string) (string, error)

	//Set sets a string in cache using provided key. Returns error if setting fails.
	Set(key string, value string, d time.Duration) error
}

type MemoryCache struct {
	cache *cache.Cache
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: cache.New(time.Hour*2, time.Minute*10),
	}
}

func (c *MemoryCache) Get(key string) (string, error) {
	result, ok := c.cache.Get(key)
	if !ok {
		return "", errors.New("no result in cache")
	}

	return result.(string), nil
}

func (c *MemoryCache) Set(key string, value string, d time.Duration) error {
	c.cache.Set(key, value, d)
	return nil
}
