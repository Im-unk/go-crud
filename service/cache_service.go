package service

import (
	"fmt"
	"time"

	"main.go/cache"
)

type CacheService struct {
	cache cache.Cacher
}

func NewCacheService(cache cache.Cacher) *CacheService {
	return &CacheService{
		cache: cache,
	}
}

func (cs *CacheService) Get(key string, v interface{}) error {
	err := cs.cache.Get(key, v)
	if err != nil {
		fmt.Printf("Cache miss: %v\n", err)
		return err
	}

	fmt.Println("Cache hit")
	return nil
}

func (cs *CacheService) Set(key string, v interface{}, expiration time.Duration) error {
	err := cs.cache.Set(key, v, expiration)
	if err != nil {
		fmt.Printf("Failed to set cache: %v\n", err)
		return err
	}

	return nil
}

func (cs *CacheService) Delete(key string) error {
	err := cs.cache.Delete(key)
	if err != nil {
		fmt.Printf("Failed to delete cache: %v\n", err)
		return err
	}

	return nil
}
