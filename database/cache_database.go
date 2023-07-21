package database

import (
	"fmt"
	"time"

	"main.go/cache"
)

type CacheDatabase struct {
	cache cache.Cacher
}

func NewCacheDatabase(cache cache.Cacher) *CacheDatabase {
	return &CacheDatabase{
		cache: cache,
	}
}

func (cs *CacheDatabase) Get(key string, v interface{}) error {
	err := cs.cache.Get(key, v)
	if err != nil {
		fmt.Printf("Cache miss: %v\n", err)
		return err
	}

	fmt.Println("Cache hit")
	return nil
}

func (cs *CacheDatabase) Set(key string, v interface{}, expiration time.Duration) error {
	err := cs.cache.Set(key, v, expiration)
	if err != nil {
		fmt.Printf("Failed to set cache: %v\n", err)
		return err
	}

	return nil
}

func (cs *CacheDatabase) Delete(key string) error {
	err := cs.cache.Delete(key)
	if err != nil {
		fmt.Printf("Failed to delete cache: %v\n", err)
		return err
	}

	return nil
}
