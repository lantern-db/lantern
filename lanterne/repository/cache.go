package repository

import (
	"context"
	"github.com/piroyoung/lanterne/lanterne/cache"
	"github.com/piroyoung/lanterne/lanterne/model"
	"sync"
	"time"
)

type CacheGraphRepository struct {
	vertices cache.VertexCache
	edges    cache.EdgeCache
}

func NewCacheGraphRepository(ttl time.Duration) CacheGraphRepository {
	return CacheGraphRepository{
		vertices: cache.NewVertexCache(ttl),
		edges:    cache.NewEdgeCache(ttl),
	}
}

func (c *CacheGraphRepository) LoadNeighbor(ctx context.Context, query model.NeighborQuery) (model.Graph, error) {
	g := model.NewGraph()
	g.VertexMap[query.Seed.Digest()] = query.Seed
	seen := map[string]model.Vertex{}

	for i := 0; i <= query.Degree - 1; i++ {
		g, seen = c.expand(query, g, seen)
	}

	return g, nil
}

func (c *CacheGraphRepository) DumpVertex(ctx context.Context, vertex model.Vertex) error {
	c.vertices.Set(vertex.Digest(), vertex)
	return nil
}

func (c *CacheGraphRepository) DumpEdge(ctx context.Context, edge model.Edge) error {
	c.vertices.Set(edge.Tail.Digest(), edge.Tail)
	c.vertices.Set(edge.Head.Digest(), edge.Head)
	c.edges.Set(edge.Tail.Digest(), edge.Head.Digest(), edge.Weight)

	return nil
}

func (c *CacheGraphRepository) getAdjacentGraph(query model.NeighborQuery, tail model.Vertex, ch chan model.Graph, wg *sync.WaitGroup) {
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

func (c *CacheGraphRepository) expand(query model.NeighborQuery, graph model.Graph, seen map[string]model.Vertex) (model.Graph, map[string]model.Vertex) {
	var wg sync.WaitGroup
	ch := make(chan model.Graph)

	nextSeen := make(map[string]model.Vertex)
	for k, v := range seen {
		nextSeen[k] = v
	}

	for digest, vertex := range graph.VertexMap {
		_, ok := seen[digest]
		if ok {
			continue
		}
		wg.Add(1)
		nextSeen[vertex.Digest()] = vertex
		go c.getAdjacentGraph(query, vertex, ch, &wg)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- graph
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := model.NewGraph()
	for graph := range ch {
		for k, v := range graph.VertexMap {
			result.VertexMap[k] = v
		}
		for tail, headMap := range graph.EdgeMap {
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
