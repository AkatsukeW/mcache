package cache

import (
	"sync"
)

type inMemoryCache struct {
	c map[string][]byte
	sync.RWMutex
	Stat
}

func (c *inMemoryCache) Set(key string, value []byte) error {
	c.Lock()
	defer c.Unlock()

	tmp, exist := c.c[key]
	if exist {
		c.del(key, tmp)
	}

	c.c[key] = value
	c.add(key, value)
	return nil
}

func (c *inMemoryCache) Get(key string) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	return c.c[key], nil
}

func (c *inMemoryCache) Del(key string) error {
	c.Lock()
	defer c.Unlock()

	v, exist := c.c[key]
	if exist {
		delete(c.c, key)
		c.del(key, v)
	}

	return nil
}

func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		make(map[string][]byte),
		sync.RWMutex{},
		Stat{},
	}
}
