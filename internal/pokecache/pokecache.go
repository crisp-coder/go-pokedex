package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache_map map[string]cacheEntry
	mu        sync.Mutex
}

type cacheEntry struct {
	key        string
	bytes      []byte
	created_at time.Time
}

func NewCache(t time.Duration) *Cache {
	c := &Cache{
		cache_map: make(map[string]cacheEntry, 0),
		mu:        sync.Mutex{},
	}
	go c.reapLoop(t)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.cache_map[key] = cacheEntry{key: key, bytes: val, created_at: time.Now()}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (cacheEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if val, ok := c.cache_map[key]; ok {
		return val, true
	} else {
		return cacheEntry{}, false
	}
}

func (c *Cache) reapLoop(t time.Duration) {
	if t > 0 {
		reapTicker := time.NewTicker(t)
		for tick := range reapTicker.C {
			c.mu.Lock()
			for key, entry := range c.cache_map {
				if tick.Sub(entry.created_at) > t {
					delete(c.cache_map, key)
				}
			}
			c.mu.Unlock()
		}
	}
}
