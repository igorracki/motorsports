package cache

import (
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any, duration time.Duration)
	Delete(key string)
	Close()
}

type cacheItem struct {
	value      any
	expiration int64
}

type memoryCache struct {
	items  map[string]cacheItem
	mu     sync.RWMutex
	stopCh chan struct{}
}

func NewMemoryCache() Cache {
	cache := &memoryCache{
		items:  make(map[string]cacheItem),
		stopCh: make(chan struct{}),
	}
	go cache.cleanup()
	return cache
}

func (cache *memoryCache) Get(key string) (any, bool) {
	cache.mu.RLock()
	item, found := cache.items[key]

	if !found {
		cache.mu.RUnlock()
		return nil, false
	}

	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		cache.mu.RUnlock()
		cache.mu.Lock()
		defer cache.mu.Unlock()

		// Re-verify after getting write lock
		if item, found = cache.items[key]; found && item.expiration > 0 && time.Now().UnixNano() > item.expiration {
			delete(cache.items, key)
		}
		return nil, false
	}

	value := item.value
	cache.mu.RUnlock()
	return value, true
}

func (cache *memoryCache) Set(key string, value any, duration time.Duration) {
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.items[key] = cacheItem{
		value:      value,
		expiration: expiration,
	}
}

func (cache *memoryCache) Delete(key string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	delete(cache.items, key)
}

func (cache *memoryCache) Close() {
	close(cache.stopCh)
}

func (cache *memoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cache.mu.Lock()
			now := time.Now().UnixNano()
			for key, item := range cache.items {
				if item.expiration > 0 && now > item.expiration {
					delete(cache.items, key)
				}
			}
			cache.mu.Unlock()
		case <-cache.stopCh:
			return
		}
	}
}
