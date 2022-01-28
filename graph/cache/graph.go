package cache

import (
	"github.com/lantern-db/lantern/graph/model"
	"sync"
)

type ValueVertex struct {
	K string
	V interface{}
}

func (v *ValueVertex) Key() string {
	return v.K
}

func (v *ValueVertex) Value() interface{} {
	return v.V
}

type GraphCache struct {
	vertices *VertexCache
	edges    *EdgeCache
}

func NewGraphCache(vertexCache *VertexCache, edgeCache *EdgeCache) *GraphCache {
	return &GraphCache{
		vertices: vertexCache,
		edges:    edgeCache,
	}
}

func (c *GraphCache) Load(query model.LoadQuery) model.Graph {
	g := model.NewGraph()
	loadedSeed, found := c.vertices.Get(query.Seed.Key())
	if found {
		g.VertexMap[query.Seed.Key()] = loadedSeed

	} else {
		g.VertexMap[query.Seed.Key()] = &ValueVertex{
			K: query.Seed.Key(),
			V: nil,
		}
	}

	seen := map[string]model.Vertex{}

	for i := uint32(0); i < query.Step; i++ {
		g, seen = c.expand(query, g, seen)
	}

	return g
}

func (c *GraphCache) LoadVertex(key string) (model.Vertex, bool) {
	return c.vertices.Get(key)
}

func (c *GraphCache) DumpVertex(vertex model.Vertex) {
	c.vertices.Set(vertex.Key(), vertex)
}

func (c *GraphCache) DumpEdge(edge model.Edge) {
	c.vertices.Set(edge.Tail.Key(), edge.Tail)
	c.vertices.Set(edge.Head.Key(), edge.Head)
	c.edges.Set(edge.Tail.Key(), edge.Head.Key(), edge.Weight)
}

func (c *GraphCache) calculateAdjacent(query model.LoadQuery, tail model.Vertex, ch chan model.Graph, wg *sync.WaitGroup) {
	defer wg.Done()

	result := model.NewGraph()
	result.VertexMap[tail.Key()] = tail
	heads, found := c.edges.GetAdjacent(tail.Key())
	if !found {
		ch <- result
		return
	}
	for headDigest, weight := range heads {
		if query.MinWeight <= weight && weight <= query.MaxWeight {
			head, found := c.vertices.Get(headDigest)
			if found {
				_, ok := result.EdgeMap[tail.Key()]
				if !ok {
					result.EdgeMap[tail.Key()] = make(map[string]float32)
				}
				result.VertexMap[head.Key()] = head
				result.EdgeMap[tail.Key()][head.Key()] = weight
			}
		}
	}

	ch <- result
	return
}

func (c *GraphCache) expand(query model.LoadQuery, graph model.Graph, seen map[string]model.Vertex) (model.Graph, map[string]model.Vertex) {
	var wg sync.WaitGroup
	ch := make(chan model.Graph)

	nextSeen := make(map[string]model.Vertex)
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
		nextSeen[vertex.Key()] = vertex
		wg.Add(1)
		go c.calculateAdjacent(query, vertex, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := model.NewGraph()
	for g := range ch {
		for k, v := range g.VertexMap {
			result.VertexMap[k] = v
		}
		for tail, headMap := range g.EdgeMap {
			for head, weight := range headMap {
				_, ok := result.EdgeMap[tail]
				if !ok {
					result.EdgeMap[tail] = make(map[string]float32)
				}
				result.EdgeMap[tail][head] = weight
			}
		}
	}

	return result, nextSeen
}

func (c *GraphCache) Flush() {
	c.vertices.Flush()
	c.edges.Flush()
}
