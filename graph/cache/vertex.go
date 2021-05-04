package cache

import (
	"github.com/piroyoung/lanterne/graph/model"
	"sync"
	"time"
)

type itemVertex struct {
	value      model.Vertex
	expiration int64
}

type VertexCache struct {
	ttl   time.Duration
	cache map[string]*itemVertex
	mu    sync.RWMutex
}

func NewVertexCache(ttl time.Duration) VertexCache {
	return VertexCache{
		ttl:   ttl,
		cache: make(map[string]*itemVertex),
	}
}

func (c *VertexCache) Set(key string, vertex model.Vertex) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = &itemVertex{
		value:      vertex,
		expiration: time.Now().Add(c.ttl).Unix(),
	}
}

func (c *VertexCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.cache[key]
	if ok {
		delete(c.cache, key)
	}
}

func (c *VertexCache) Get(key string) (model.Vertex, bool) {
	c.mu.RLock()

	item, ok := c.cache[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}
	if time.Now().Unix() > item.expiration {
		defer c.Delete(key)
		return nil, false
	}
	return item.value, true
}

func (c *VertexCache) Flush() {
	var keys []string
	c.mu.RLock()
	for key, vertex := range c.cache {
		if time.Now().Unix() > vertex.expiration {
			keys = append(keys, key)
		}
	}
	c.mu.RUnlock()

	for _, key := range keys {
		c.Delete(key)
	}
}
