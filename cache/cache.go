package cache

import (
	"time"
)

type Cacher interface {
	Get(key string, v interface{}) error
	Set(key string, v interface{}, expiration time.Duration) error
	Delete(key string) error
}
