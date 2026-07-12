package cache

import (
	"slices"
	"sync"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func New() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(key string, value []byte) {
	v := slices.Clone(value)

	c.mu.Lock()
	c.data[key] = v
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	value, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	return append([]byte(nil), value...), true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

func (c *Cache) Has(key string) bool {
	c.mu.RLock()
	_, ok := c.data[key]
	c.mu.RUnlock()
	return ok
}

func (c *Cache) Len() int {
	c.mu.RLock()
	n := len(c.data)
	c.mu.RUnlock()
	return n
}

func (c *Cache) Flush() {
	c.mu.Lock()
	c.data = make(map[string][]byte)
	c.mu.Unlock()
}
