package cache

import (
	"github.com/piroyoung/lanterne/lanterne/model"
	"sync"
	"time"
)

type Item struct {
	value      model.Vertex
	expiration int64
}

type VertexCache struct {
	ttl   time.Duration
	cache map[string]*Item
	mu    sync.RWMutex
}

func NewVertexCache(ttl time.Duration) VertexCache {
	return VertexCache{
		ttl:   ttl,
		cache: make(map[string]*Item),
	}
}

func (c *VertexCache) Set(digest string, vertex model.Vertex) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[digest] = &Item{
		value:      vertex,
		expiration: time.Now().Add(c.ttl).Unix(),
	}
}

func (c *VertexCache) Delete(digest string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.cache[digest]
	if ok {
		delete(c.cache, digest)
	}
}

func (c *VertexCache) Get(digest string) (model.Vertex, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.cache[digest]
	if !ok {
		return nil, false
	}
	if time.Now().Unix() > item.expiration {
		c.Delete(digest)
		return nil, false
	}
	return item.value, true
}
