package pokecache

import (
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	Entry map[string]cacheEntry
	mutex    sync.Mutex
}

func (c *Cache) reapLoop(interval time.Duration) {
    ticker := time.NewTicker(interval)
    
    go func() {
        for {
            <-ticker.C
            
            c.mutex.Lock()
            
            now := time.Now()
            for key, Entry := range c.Entry {
                if now.Sub(Entry.createdAt) > interval {
                    delete(c.Entry, key)
                }
            }
            
            c.mutex.Unlock()
        }
    }()
}


func NewCache(interval time.Duration) *Cache {
	cachePtr := &Cache{
		Entry: make(map[string]cacheEntry),
	}

	cachePtr.reapLoop(interval)


	return cachePtr
}

func (c *Cache) Add(key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Entry[key] = cacheEntry{
		createdAt: time.Now(),
		val: value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	val, exists := c.Entry[key]
	if exists {
		return val.val, true
	} else {
		return []byte{}, false
	}
}