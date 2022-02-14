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

	c.cache[vertex.Key()] = vertex
}

func (c *VertexCache) delete(key Key) {
	if _, ok := c.cache[key]; ok {
		delete(c.cache, key)
	}
}

func (c *VertexCache) Delete(key Key) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.delete(key)
}

func (c *VertexCache) Get(key Key) (Vertex, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.cache[key]; !ok {
		return nil, false

	} else if item.Expiration().Dead() {
		go c.delete(key)
		return nil, false

	} else {
		return item, true

	}
}

func (c *VertexCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, vertex := range c.cache {
		if vertex.Expiration().Dead() {
			c.delete(key)
		}
	}
}
