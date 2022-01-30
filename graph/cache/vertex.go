package cache

import (
	. "github.com/lantern-db/lantern/graph/model"
	"sync"
)

type VertexCache struct {
	cache map[Key]Vertex
	mu    sync.RWMutex
}

func NewVertexCache() *VertexCache {
	return &VertexCache{
		cache: make(map[Key]Vertex),
	}
}

func (c *VertexCache) Set(vertex Vertex) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[vertex.Key] = vertex
}

func (c *VertexCache) Delete(key Key) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.cache[key]
	if ok {
		delete(c.cache, key)
	}
}

func (c *VertexCache) Get(key Key) (Vertex, bool) {
	c.mu.RLock()

	item, ok := c.cache[key]
	c.mu.RUnlock()

	if !ok {
		return Vertex{}, false
	}
	if item.Expiration.Dead() {
		defer c.Delete(key)
		return Vertex{}, false
	}
	return item, true
}

func (c *VertexCache) Has(key Key) bool {
	_, ok := c.cache[key]
	return ok
}

func (c *VertexCache) Flush() {
	c.mu.RLock()
	for key, vertex := range c.cache {
		if vertex.Expiration.Dead() {
			c.Delete(key)
		}
	}
	c.mu.RUnlock()
}
