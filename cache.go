package ttlcache

import (
	"sync"
	"time"
)

// Cache is a synchronised map of items that auto-expire once stale
type Cache struct {
	mutex sync.RWMutex
	ttl   time.Duration
	items map[string]*Item
	isTTL bool
}

// Set is a thread-safe way to add new items to the map
func (cache *Cache) Set(key string, data interface{}) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item := &Item{data: data}
	if cache.isTTL {
		item.touch(cache.ttl)
	}
	cache.items[key] = item
}

// SetTTL is a thread-safe way to add new items to the map with time TTL
func (cache *Cache) SetTTL(key string, data interface{}, ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if !cache.isTTL {
		cache.isTTL = true
		cache.startTTLCleanupTimer()
	}
	item := &Item{data: data}
	item.touch(ttl)
	cache.items[key] = item
}

// Get is a thread-safe way to lookup items
// Every lookup, also touches the item, hence extending it's life
func (cache *Cache) Get(key string) (data interface{}, found bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	item, exists := cache.items[key]
	if cache.isTTL {
		if !exists || item.expired() {
			data = nil
			found = false
		} else {
			item.touch(cache.ttl)
			data = item.data
			found = true
		}
	} else {
		if exists {
			return item.data, true
		}
		return nil, false
	}
	return
}

// Count returns the number of items in the cache
// (helpful for tracking memory leaks)
func (cache *Cache) Count() int {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	count := len(cache.items)
	return count
}

// Clear removes all entries from the cache
func (cache *Cache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	for key := range cache.items {
		delete(cache.items, key)
	}
}

func (cache *Cache) ttlCleanup() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if cache.isTTL {
		for key, item := range cache.items {
			if item.expired() {
				delete(cache.items, key)
			}
		}
	}
}

// Delete removes an entry from the cache at the specified key.
// If no entry exists at the specified key, no action is taken
func (cache *Cache) Delete(key string) {
	if _, ok := cache.items[key]; ok {
		delete(cache.items, key)
	}
}

// startTTLCleanupTimer installs a timer to perform
// cache entry removal if item has reached time to live.
func (cache *Cache) startTTLCleanupTimer() {
	if cache.isTTL {
		duration := cache.ttl
		if duration < time.Millisecond {
			duration = time.Millisecond
		}
		ticker := time.Tick(duration)
		go (func() {
			for {
				select {
				case <-ticker:
					cache.ttlCleanup()
				}
			}
		})()
	}
}

// NewTTLCache is a helper to create instance of the Cache struct
func NewTTLCache(duration time.Duration) *Cache {
	cache := &Cache{
		ttl:   duration,
		items: map[string]*Item{},
		isTTL: true,
	}
	cache.startTTLCleanupTimer()
	return cache
}

// NewCache is a helper to create instance of the Cache struct
func NewCache() *Cache {
	cache := &Cache{
		items: map[string]*Item{},
		isTTL: false,
	}
	return cache
}
