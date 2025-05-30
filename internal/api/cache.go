package api

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	item     map[string]cacheEntry
	mu       sync.RWMutex
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		item:     make(map[string]cacheEntry),
		interval: interval,
	}

	go cache.ReapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.item[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.item[key]
	if !ok {
		return nil, false
	}

	return item.val, true
}

func (c *Cache) ReapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Reap()
		}
	}
}

func (c *Cache) Reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	expiredKeys := make([]string, 0)
	for key, item := range c.item {
		if now.Sub(item.createdAt) > c.interval {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		delete(c.item, key)
	}

	if len(expiredKeys) > 0 {
		fmt.Printf("Reaped %d expired cache entries\n", len(expiredKeys))
	}
}
