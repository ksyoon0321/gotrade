package util

import (
	"sync"
	"time"
)

type CacheItem struct {
	v      interface{}
	expire *time.Time
}
type CacheTimer struct {
	ttl   time.Duration
	items map[string]*CacheItem
	lock  *sync.RWMutex
}

func NewCacheTimer(interval time.Duration) *CacheTimer {
	if interval < time.Second {
		interval = time.Second
	}

	cache := &CacheTimer{
		ttl:   interval,
		items: make(map[string]*CacheItem),
		lock:  &sync.RWMutex{},
	}

	go func() {
		tick := time.NewTicker(cache.ttl)
		for {
			now := <-tick.C

			cache.lock.Lock()
			for id, item := range cache.items {
				if item.expire != nil && item.expire.Before(now) {
					delete(cache.items, id)
				}
			}
			cache.lock.Unlock()
		}
	}()
	return cache
}

func (c *CacheTimer) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, ok := c.items[key]

	if ok && item.expire != nil && item.expire.After(time.Now()) {
		return item.v
	}
	return nil
}

func (c *CacheTimer) Put(key string, v interface{}, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	expire := time.Now().Add(ttl)

	c.items[key] = &CacheItem{
		v:      v,
		expire: &expire,
	}
}
