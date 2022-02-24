package cache

import (
	. "github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/graph/table"
	"sync"
)

type EdgeCache struct {
	cache          map[Key]map[Key]*table.EdgeTable
	incomingDegree IncrementalStatMap
	outgoingDegree IncrementalStatMap
	mu             sync.RWMutex
}

func NewEdgeCache() *EdgeCache {
	return &EdgeCache{
		cache:          make(map[Key]map[Key]*table.EdgeTable),
		incomingDegree: make(IncrementalStatMap),
		outgoingDegree: make(IncrementalStatMap),
	}
}

func (c *EdgeCache) Put(edge Edge) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.cache[edge.Tail()]; !ok {
		c.cache[edge.Tail()] = make(map[Key]*table.EdgeTable)
	}

	if _, ok := c.cache[edge.Tail()][edge.Head()]; !ok {
		c.cache[edge.Tail()][edge.Head()] = table.NewEmptyEdgeTable()
		c.incomingDegree.Increment(edge.Head())
		c.outgoingDegree.Increment(edge.Tail())
	}

	c.cache[edge.Tail()][edge.Head()].Append(edge)
}

func (c *EdgeCache) delete(tail Key, head Key) {
	if _, ok := c.cache[tail]; ok {
		if _, ok := c.cache[tail][head]; ok {
			delete(c.cache[tail], head)
			c.outgoingDegree.Decrement(tail)
			if len(c.cache[tail]) == 0 {
				delete(c.cache, tail)
				c.incomingDegree.Decrement(head)
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

	if edgeTable, ok := c.cache[tail][head]; !ok {
		return nil, false

	} else if edgeTable.IsEmpty() {
		go c.delete(tail, head)
		return nil, false

	} else {
		edge := NewStaticEdge(tail, head, edgeTable.Weight(), edgeTable.Expiration())
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
	for head, edgeTable := range headMap {
		if edgeTable.IsEmpty() {
			go c.delete(tail, head)

		} else {
			result[head] = NewStaticEdge(tail, head, edgeTable.Weight(), edgeTable.Expiration())

		}
	}

	return result, len(result) != 0
}

func (c *EdgeCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for tail, headMap := range c.cache {
		for head, edgeTable := range headMap {
			if edgeTable.IsEmpty() {
				c.delete(tail, head)
			}
		}
	}
}
