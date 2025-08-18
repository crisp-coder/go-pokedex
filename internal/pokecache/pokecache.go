package pokecache

import (
	"sync"
	"time"
)

type PokeCache struct {
	cache_map map[string]cacheEntry
	mu        sync.Mutex
}

type cacheEntry struct {
	bytes      []byte
	created_at time.Time
}

func NewPokeCache(t time.Duration) *PokeCache {
	c := &PokeCache{
		cache_map: make(map[string]cacheEntry, 0),
		mu:        sync.Mutex{},
	}
	go c.reapLoop(t)
	return c
}

func (c *PokeCache) Add(key string, val []byte) {
	c.mu.Lock()
	c.cache_map[key] = cacheEntry{bytes: val, created_at: time.Now()}
	c.mu.Unlock()
}

func (c *PokeCache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if val, ok := c.cache_map[key]; ok {
		return val.bytes, true
	} else {
		return []byte{}, false
	}
}

func (c *PokeCache) reapLoop(t time.Duration) {
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
