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

	_, ok := c.cache[edge.Tail()]
	if !ok {
		c.cache[edge.Tail()] = make(map[Key]Edge)
	}
	c.cache[edge.Tail()][edge.Head()] = edge
}

func (c *EdgeCache) delete(tail Key, head Key) {
	if _, ok := c.cache[tail]; ok {
		if _, ok := c.cache[tail][head]; ok {
			delete(c.cache[tail], head)
			if len(c.cache[tail]) == 0 {
				delete(c.cache, tail)
			}
		}
	}
}

func (c *EdgeCache) Delete(tail Key, head Key) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.delete(tail, head)
}

func (c *EdgeCache) Get(tail Key, head Key) (Edge, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.cache[tail]; !ok {
		return nil, false
	}

	if edge, ok := c.cache[tail][head]; !ok {
		return nil, false

	} else if edge.Expiration().Dead() {
		go c.delete(tail, head)
		return nil, false

	} else {
		return edge, true

	}
}

func (c *EdgeCache) GetAdjacent(tail Key) (map[Key]Edge, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	headMap, ok := c.cache[tail]
	if !ok {
		return nil, false
	}

	result := make(map[Key]Edge)
	for head, edge := range headMap {
		if edge.Expiration().Dead() {
			go c.delete(tail, head)

		} else {
			result[head] = edge

		}
	}

	return result, len(result) != 0
}

func (c *EdgeCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for tail, headMap := range c.cache {
		for head, edge := range headMap {
			if edge.Expiration().Dead() {
				c.delete(tail, head)
			}
		}
	}
}
