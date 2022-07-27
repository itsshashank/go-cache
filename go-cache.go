package gocache

import (
	"fmt"
	"sync"
	"time"
)

type element struct {
	Object     any
	expiration int64
}

type Cache struct {
	elements map[any]element
	mu       sync.RWMutex
}

func (c *Cache) deleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.elements {
		if now > v.expiration {
			delete(c.elements, k)
		}
	}
	c.mu.Unlock()
}

type watchdog struct {
	interval time.Duration
}

func (w *watchdog) Run(c *Cache) {
	ticker := time.NewTicker(w.interval)
	for range ticker.C {
		c.deleteExpired()
	}
}

func runWatchdog(c *Cache, intervel time.Duration) {
	wd := watchdog{interval: intervel}
	go wd.Run(c)
}

func New() *Cache {
	cache := &Cache{elements: make(map[any]element)}
	runWatchdog(cache, 10) //default intervel of 10 sec
	return cache
}

func (c *Cache) Get(key any) (any, bool) {
	defer c.mu.RUnlock()
	c.mu.RLock()
	element, found := c.elements[key]
	if !found {
		return nil, found
	}
	if element.expiration > 0 {
		if time.Now().UnixNano() > element.expiration {
			return nil, false
		}
	}
	return element.Object, true
}

func (c *Cache) Set(key any, value any, duriation time.Duration) {
	if len(c.elements) >= 1024 {
		fmt.Printf("Maximum Cache elements present cant set { %v, %v }", key, value)
		return
	}
	var endtime int64
	if duriation > 0 {
		endtime = time.Now().Add(duriation).UnixNano()
	}
	c.mu.Lock()
	c.elements[key] = element{
		value,
		endtime,
	}
	c.mu.Unlock()
}
