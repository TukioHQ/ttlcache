## TTLCache - an in-memory LRU cache with expiration

TTLCache is a minimal wrapper over a interface type map in golang, entries of which are

1. Thread-safe
2. Auto-Expiring after a certain time
3. Auto-Extending expiration on `Get`s

<!-- [![Build Status](https://travis-ci.org/TukioHQ/ttlcache.svg)](https://travis-ci.org/TukioHQ/ttlcache) -->

#### TTL Usage
```go
import (
  "time"
  "github.com/TukioHQ/ttlcache"
)

func main () {
  cache := ttlcache.NewTTLCache(time.Second)
  cache.Set("key", 200)
  value, exists := cache.Get("key")
  count := cache.Count()
}
```

#### Non TTL Usage
```go
import (
  "time"
  "github.com/TukioHQ/ttlcache"
)

func main () {
  cache := ttlcache.NewCache()
  cache.Set("key", 200)
  value, exists := cache.Get("key")
  count := cache.Count()
}
```