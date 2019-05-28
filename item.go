package ttlcache

import (
	"sync"
	"time"
)

// Item represents a record in the cache map
type Item struct {
	sync.RWMutex
	data    interface{}
	expires *time.Time
	isTTL   bool
}

func (item *Item) touch(duration time.Duration) {
	item.Lock()
	defer item.Unlock()
	expiration := time.Now().Add(duration)
	item.expires = &expiration
	item.isTTL = true
}

func (item *Item) expired() bool {
	item.RLock()
	defer item.RUnlock()
	if !item.isTTL {
		return false
	}
	var value bool
	if item.expires == nil {
		value = true
	} else {
		value = item.expires.Before(time.Now())
	}
	return value
}
