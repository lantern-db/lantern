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

	heads, found := c.edges.GetHeads(tail.Digest())
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
	c.vertices.Set(vertex.Digest(), vertex)
	return nil
}

func (c *CacheGraphRepository) DumpEdge(ctx context.Context, edge model.Edge) error {
	if err := c.DumpVertex(ctx, edge.Tail); err != nil {
		return err
	}
	if err := c.DumpVertex(ctx, edge.Head); err != nil {
		return err
	}
	c.edges.Set(edge.Tail.Digest(), edge.Head.Digest(), edge.Weight)

	return nil
}
