package cache

import (
	"github.com/lantern-db/lantern/graph/model"
	"sync"
)

type EdgeCache struct {
	cache map[model.Key]map[model.Key]model.Edge
	mu    sync.RWMutex
}

func NewEdgeCache() *EdgeCache {
	return &EdgeCache{
		cache: make(map[model.Key]map[model.Key]model.Edge),
	}
}

func (c *EdgeCache) Set(edge model.Edge) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.cache[edge.Tail]
	if !ok {
		c.cache[edge.Tail] = make(map[model.Key]model.Edge)
	}
	c.cache[edge.Tail][edge.Head] = edge
}

func (c *EdgeCache) Delete(tail model.Key, head model.Key) {
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

func (c *EdgeCache) Get(tail model.Key, head model.Key) (model.Edge, bool) {
	c.mu.RLock()

	if _, ok := c.cache[tail]; !ok {
		c.mu.RUnlock()
		return model.Edge{}, false
	}

	edge, ok := c.cache[tail][head]
	c.mu.RUnlock()

	if !ok {
		return model.Edge{}, false
	}

	if edge.Expiration.Dead() {
		c.Delete(tail, head)
		return model.Edge{}, false
	}
	return edge, true
}

func (c *EdgeCache) GetAdjacent(tail model.Key) (map[model.Key]model.Edge, bool) {
	result := make(map[model.Key]model.Edge)
	var expired []model.Key

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
	var keys []model.Edge
	c.mu.RLock()
	for tail, headMap := range c.cache {
		for head, edge := range headMap {
			if edge.Expiration.Dead() {
				keys = append(keys, model.Edge{Tail: tail, Head: head})
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
