package cache

import (
	. "github.com/lantern-db/lantern/graph/model"
	"sync"
)

type EdgeCache struct {
	cache map[Key]map[Key]Edge
	mu    sync.RWMutex
}

func NewEdgeCache() *EdgeCache {
	return &EdgeCache{
		cache: make(map[Key]map[Key]Edge),
	}
}

func (c *EdgeCache) Set(edge Edge) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.cache[edge.Tail]
	if !ok {
		c.cache[edge.Tail] = make(map[Key]Edge)
	}
	c.cache[edge.Tail][edge.Head] = edge
}

func (c *EdgeCache) Delete(tail Key, head Key) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.cache[tail]; ok {
		if _, ok := c.cache[tail][head]; ok {
			delete(c.cache[tail], head)
			if len(c.cache[tail]) == 0 {
				delete(c.cache, tail)
			}
		}
	}
}

func (c *EdgeCache) Get(tail Key, head Key) (Edge, bool) {
	c.mu.RLock()

	if _, ok := c.cache[tail]; !ok {
		c.mu.RUnlock()
		return Edge{}, false
	}

	edge, ok := c.cache[tail][head]
	c.mu.RUnlock()

	if !ok {
		return Edge{}, false
	}

	if edge.Expiration.Dead() {
		c.Delete(tail, head)
		return Edge{}, false
	}
	return edge, true
}

func (c *EdgeCache) GetAdjacent(tail Key) (map[Key]Edge, bool) {
	result := make(map[Key]Edge)
	var expired []Key

	c.mu.RLock()
	headMap, ok := c.cache[tail]
	if !ok {
		return nil, false
	}
	for head, edge := range headMap {
		if edge.Expiration.Dead() {
			expired = append(expired, head)
		} else {
			result[head] = edge
		}
	}
	c.mu.RUnlock()

	go func() {
		for _, head := range expired {
			c.Delete(tail, head)
		}
	}()

	return result, len(result) != 0
}

func (c *EdgeCache) Flush() {
	var keys []Edge
	c.mu.RLock()
	for tail, headMap := range c.cache {
		for head, edge := range headMap {
			if edge.Expiration.Dead() {
				keys = append(keys, Edge{Tail: tail, Head: head})
			}
		}
	}
	c.mu.RUnlock()
	go func() {
		for _, key := range keys {
			c.Delete(key.Tail, key.Head)
		}
	}()
}
