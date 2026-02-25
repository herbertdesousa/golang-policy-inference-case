package cache

import (
	"log"

	lru "github.com/hashicorp/golang-lru/v2"
)

type Cache[T any] struct {
	cache *lru.Cache[string, *T]
}

func NewCache[T any](size int) (*Cache[T], error) {
	cache, err := lru.New[string, *T](size)

	if err != nil {
		log.Fatalf("Failed to create cache: %v", err)
		return nil, err
	}

	return &Cache[T]{cache: cache}, nil
}

func (c *Cache[T]) Get(key string) (*T, bool) {
	return c.cache.Get(key)
}

func (c *Cache[T]) Add(key string, value *T) {
	c.cache.Add(key, value)
}
