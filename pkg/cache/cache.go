package cache

import (
	"strings"
	"sync"
	"time"
)

type cacher struct {
	store    map[string]entry
	rwMutex  sync.RWMutex
	lifespan time.Duration
}

type entry struct {
	value   any
	created time.Time
}

var Cache *cacher

func Init(lifespan time.Duration) {
	if lifespan == 0 {
		lifespan = time.Second * 15
	}

	Cache = &cacher{
		store:    make(map[string]entry),
		lifespan: lifespan,
	}
}

func (c *cacher) Get(key string) (value any, ok bool) {
	c.rwMutex.RLock()
	e, ok := c.store[key]
	c.rwMutex.RUnlock()

	if ok {
		if e.created.Before(time.Now().Add(-c.lifespan)) {
			c.rwMutex.Lock()
			delete(c.store, key)
			c.rwMutex.Unlock()

			return nil, false
		}
	}

	value = e.value
	return
}

func (c *cacher) Set(key string, value any) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.store[key] = entry{
		value:   value,
		created: time.Now(),
	}
}

func (c *cacher) Invalidate(partOfKey string) {
	c.rwMutex.RLock()
	var toInvalidate []string

	for k, _ := range c.store {
		if strings.Contains(k, partOfKey) {
			toInvalidate = append(toInvalidate, k)
		}
	}
	c.rwMutex.RUnlock()

	for _, k := range toInvalidate {
		c.rwMutex.Lock()
		delete(c.store, k)
		c.rwMutex.Unlock()
	}
}
