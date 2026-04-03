package cache

import (
	"sync"
	"time"
)

type Cache[T any] interface {
	Get(key string) (T, bool)
	Set(key string, value T, duration time.Duration)
	Delete(key string)
	Close()
}

type cacheItem[T any] struct {
	value      T
	expiration int64
}

type memoryCache[T any] struct {
	items       map[string]cacheItem[T]
	mutex       sync.RWMutex
	stopChannel chan struct{}
}

func NewMemoryCache[T any]() Cache[T] {
	cache := &memoryCache[T]{
		items:       make(map[string]cacheItem[T]),
		stopChannel: make(chan struct{}),
	}
	go cache.cleanup()
	return cache
}

func (cache *memoryCache[T]) Get(key string) (T, bool) {
	cache.mutex.RLock()
	item, found := cache.items[key]

	if !found {
		cache.mutex.RUnlock()
		var zero T
		return zero, false
	}

	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		cache.mutex.RUnlock()
		cache.mutex.Lock()
		defer cache.mutex.Unlock()

		if item, found = cache.items[key]; found && item.expiration > 0 && time.Now().UnixNano() > item.expiration {
			delete(cache.items, key)
		}
		var zero T
		return zero, false
	}

	value := item.value
	cache.mutex.RUnlock()
	return value, true
}

func (cache *memoryCache[T]) Set(key string, value T, duration time.Duration) {
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.items[key] = cacheItem[T]{
		value:      value,
		expiration: expiration,
	}
}

func (cache *memoryCache[T]) Delete(key string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	delete(cache.items, key)
}

func (cache *memoryCache[T]) Close() {
	close(cache.stopChannel)
}

func (cache *memoryCache[T]) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cache.mutex.Lock()
			now := time.Now().UnixNano()
			for key, item := range cache.items {
				if item.expiration > 0 && now > item.expiration {
					delete(cache.items, key)
				}
			}
			cache.mutex.Unlock()
		case <-cache.stopChannel:
			return
		}
	}
}
