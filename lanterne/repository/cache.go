package repository

import (
	"context"
	"github.com/piroyoung/lanterne/lanterne/cache"
	"github.com/piroyoung/lanterne/lanterne/model"
)

type CacheGraphRepository struct {
	vertices cache.VertexCache
	edges    cache.EdgeCache
}

func (c *CacheGraphRepository) LoadNeighbor(ctx context.Context, query model.NeighborQuery) (model.Graph, error) {
	panic("implement me")
}

func (c *CacheGraphRepository) LoadAdjacent(ctx context.Context, query model.AdjacentQuery) (model.Graph, error) {
	g := model.NewGraph()
	tail := query.Seed

	g.Vertices[tail.Digest()] = tail

	heads, found := c.edges.GetAdjacent(tail.Digest())
	if found {
		for headKey, weight := range heads {
			if query.MinWeight <= weight && weight <= query.MaxWeight {
				head, found := c.vertices.Get(headKey)
				if found {
					g.Vertices[headKey] = head
					_, ok := g.Edges[tail.Digest()]
					if !ok {
						g.Edges[tail.Digest()] = make(map[string]float32)
					}
					g.Edges[tail.Digest()][head.Digest()] = weight
				}
			}
		}
	}

	return g, nil
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

func (c *CacheGraphRepository) getAdjacentGraph(tail model.Vertex) model.Graph {
	result := model.NewGraph()
	result.Vertices[tail.Digest()] = tail
	heads, found := c.edges.GetAdjacent(tail.Digest())
	if !found {
		return result
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

	return result
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
