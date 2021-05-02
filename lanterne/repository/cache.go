package repository

import (
	"context"
	"github.com/piroyoung/lanterne/lanterne/cache"
	"github.com/piroyoung/lanterne/lanterne/model"
	"sync"
)

type CacheGraphRepository struct {
	vertices cache.VertexCache
	edges    cache.EdgeCache
}

func (c *CacheGraphRepository) LoadNeighbor(ctx context.Context, query model.NeighborQuery) (model.Graph, error) {
	panic("implement me")
}

func (c *CacheGraphRepository) DumpVertex(ctx context.Context, vertex model.Vertex) error {
	go c.vertices.Set(vertex.Digest(), vertex)
	return nil
}

func (c *CacheGraphRepository) DumpEdge(ctx context.Context, edge model.Edge) error {
	go c.vertices.Set(edge.Tail.Digest(), edge.Tail)
	go c.vertices.Set(edge.Head.Digest(), edge.Head)
	go c.edges.Set(edge.Tail.Digest(), edge.Head.Digest(), edge.Weight)

	return nil
}

func (c *CacheGraphRepository) getAdjacentGraph(tail model.Vertex, ch chan model.Graph, wg *sync.WaitGroup) {
	defer wg.Done()

	result := model.NewGraph()
	result.Vertices[tail.Digest()] = tail
	heads, found := c.edges.GetAdjacent(tail.Digest())
	if !found {
		ch <- result
		return
	}
	for headDigest, weight := range heads {
		head, found := c.vertices.Get(headDigest)
		if found {
			_, ok := result.Edges[tail.Digest()]
			if !ok {
				result.Edges[tail.Digest()] = make(map[string]float32)
			}
			result.Vertices[head.Digest()] = head
			result.Edges[tail.Digest()][head.Digest()] = weight
		}
	}

	ch <- result
	return
}

func (c *CacheGraphRepository) expand(graph model.Graph, seen map[string]model.Vertex) (model.Graph, map[string]model.Vertex) {
	var wg sync.WaitGroup
	ch := make(chan model.Graph)
	nextSeen := make(map[string]model.Vertex)
	for k, v := range seen {
		nextSeen[k] = v
	}

	for digest, vertex := range graph.Vertices {
		_, ok := seen[digest]
		if ok {
			continue
		}
		wg.Add(1)
		nextSeen[vertex.Digest()] = vertex
		go c.getAdjacentGraph(vertex, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var graphArray []model.Graph
	for g := range ch {
		graphArray = append(graphArray, g)
	}

	return union(graphArray), nextSeen
}

func union(graphArray []model.Graph) model.Graph {
	result := model.NewGraph()
	for _, graph := range graphArray {
		for k, v := range graph.Vertices {
			result.Vertices[k] = v
		}
		for tail, headMap := range result.Edges {
			for head, weight := range headMap {
				_, ok := result.Edges[tail]
				if !ok {
					result.Edges[tail] = make(map[string]float32)
				}
				result.Edges[tail][head] = weight
			}
		}
	}

	return result
}
