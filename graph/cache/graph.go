package cache

import (
	. "github.com/lantern-db/lantern/graph/model"
	"sync"
	"time"
)

type GraphCache struct {
	defaultTtl time.Duration
	vertexCache *VertexCache
	edgeCache   *EdgeCache
}

func NewGraphCache(defaultTtl time.Duration, vertexCache *VertexCache, edgeCache *EdgeCache) *GraphCache {
	return &GraphCache{
		defaultTtl: defaultTtl,
		vertexCache: vertexCache,
		edgeCache:   edgeCache,
	}
}

func (c *GraphCache) Load(query LoadQuery) Graph {
	g := NewGraph()
	loadedSeed, found := c.vertexCache.Get(query.Seed)
	if found {
		g.VertexMap[query.Seed] = loadedSeed

	} else {
		g.VertexMap[query.Seed] = Vertex{
			Key: query.Seed,
		}
	}

	seen := map[Key]Vertex{}

	for i := uint32(0); i < query.Step; i++ {
		g, seen = c.expand(query, g, seen)
	}

	return g
}

func (c *GraphCache) LoadVertex(key Key) (Vertex, bool) {
	return c.vertexCache.Get(key)
}

func (c *GraphCache) DumpVertex(vertex Vertex) {
	if vertex.Expiration == 0.0 {
		vertex.Expiration = NewExpiration(c.defaultTtl)
	}
	c.vertexCache.Set(vertex)
}

func (c *GraphCache) DumpEdge(edge Edge) {
	if edge.Expiration == 0.0 {
		edge.Expiration = NewExpiration(c.defaultTtl)
	}
	c.edgeCache.Set(edge)
}

func (c *GraphCache) calculateAdjacent(query LoadQuery, tail Vertex, ch chan Graph, wg *sync.WaitGroup) {
	defer wg.Done()

	result := NewGraph()
	result.VertexMap[tail.Key] = tail
	heads, found := c.edgeCache.GetAdjacent(tail.Key)
	if !found {
		ch <- result
		return
	}
	for headDigest, edge := range heads {
		weight := float32(edge.Weight)
		if query.MinWeight <= weight && weight <= query.MaxWeight {
			head, found := c.vertexCache.Get(headDigest)
			if found {
				_, ok := result.EdgeMap[tail.Key]
				if !ok {
					result.EdgeMap[tail.Key] = make(map[Key]Weight)
				}
				result.VertexMap[head.Key] = head
				result.EdgeMap[tail.Key][head.Key] = edge.Weight
			}
		}
	}

	ch <- result
	return
}

func (c *GraphCache) expand(query LoadQuery, graph Graph, seen map[Key]Vertex) (Graph, map[Key]Vertex) {
	var wg sync.WaitGroup
	ch := make(chan Graph)

	nextSeen := make(map[Key]Vertex)
	for k, v := range seen {
		nextSeen[k] = v
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- graph
	}()

	for digest, vertex := range graph.VertexMap {
		_, ok := seen[digest]
		if ok {
			continue
		}
		nextSeen[vertex.Key] = vertex
		wg.Add(1)
		go c.calculateAdjacent(query, vertex, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := NewGraph()
	for g := range ch {
		for k, v := range g.VertexMap {
			result.VertexMap[k] = v
		}
		for tail, headMap := range g.EdgeMap {
			for head, weight := range headMap {
				_, ok := result.EdgeMap[tail]
				if !ok {
					result.EdgeMap[tail] = make(map[Key]Weight)
				}
				result.EdgeMap[tail][head] = weight
			}
		}
	}

	return result, nextSeen
}

func (c *GraphCache) Flush() {
	c.vertexCache.Flush()
	c.edgeCache.Flush()
}
