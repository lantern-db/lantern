package cache

import (
	"github.com/piroyoung/lanterne/graph/model"
	"sync"
	"time"
)

type GraphCache struct {
	vertices VertexCache
	edges    EdgeCache
}

func NewGraphCache(ttl time.Duration) GraphCache {
	return GraphCache{
		vertices: NewVertexCache(ttl),
		edges:    NewEdgeCache(ttl),
	}
}

func (c *GraphCache) Load(query model.LoadQuery) model.Graph {
	g := model.NewGraph()
	g.VertexMap[query.Seed.Digest()] = query.Seed
	seen := map[string]model.Vertex{}

	for i := uint32(0); i+1 <= query.Degree; i++ {
		g, seen = c.expand(query, g, seen)
	}

	return g
}

func (c *GraphCache) DumpVertex(vertex model.Vertex) {
	c.vertices.Set(vertex.Digest(), vertex)
}

func (c *GraphCache) DumpEdge(edge model.Edge) {
	c.vertices.Set(edge.Tail.Digest(), edge.Tail)
	c.vertices.Set(edge.Head.Digest(), edge.Head)
	c.edges.Set(edge.Tail.Digest(), edge.Head.Digest(), edge.Weight)
}

func (c *GraphCache) calculateAdjacent(query model.LoadQuery, tail model.Vertex, ch chan model.Graph, wg *sync.WaitGroup) {
	defer wg.Done()

	result := model.NewGraph()
	result.VertexMap[tail.Digest()] = tail
	heads, found := c.edges.GetAdjacent(tail.Digest())
	if !found {
		ch <- result
		return
	}
	for headDigest, weight := range heads {
		if query.MinWeight <= weight && weight <= query.MaxWeight {
			head, found := c.vertices.Get(headDigest)
			if found {
				_, ok := result.EdgeMap[tail.Digest()]
				if !ok {
					result.EdgeMap[tail.Digest()] = make(map[string]float32)
				}
				result.VertexMap[head.Digest()] = head
				result.EdgeMap[tail.Digest()][head.Digest()] = weight
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
		nextSeen[vertex.Digest()] = vertex
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
