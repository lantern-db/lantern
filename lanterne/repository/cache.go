package repository

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/piroyoung/lanterne/lanterne/model"
	"strings"
	"time"
)

type CacheGraphRepository struct {
	client *cache.Cache
	ttl    time.Duration
}

func (c CacheGraphRepository) LoadNeighbor(ctx context.Context, query model.NeighborQuery) (model.Graph, error) {
	panic("implement me")
}

func (c CacheGraphRepository) LoadAdjacent(ctx context.Context, query model.AdjacentQuery) (model.Graph, error) {
	g := model.NewGraph()
	tail := query.Seed
	g.Vertices[tail.Digest()] = tail
	for key, value := range c.client.Items() {
		if strings.HasPrefix(key, tail.Digest()+"->") {
			switch v := value.Object.(type) {
			case float32:
				if query.MinWeight <= v && v <= query.MaxWeight {
					_, headKey := splitToHeadAndTail(key)
					head, found := c.loadVertex(ctx, headKey)
					if !found {
						continue
					}
					g.Vertices[head.Digest()] = head
					g.Edges[tail.Digest()][head.Digest()] = v
				}
			}
		}
	}
	return g, nil
}

func (c CacheGraphRepository) expandGraph(ctx context.Context, graph model.Graph, footprints map[string]bool, minWeight float32, maxWeight float32) (map[string]bool, error) {
	for _, v := range graph.Vertices {
		if _, ok := footprints[v.Digest()]; !ok {
			_, err := c.LoadAdjacent(ctx, model.AdjacentQuery{
				Seed:      v,
				MinWeight: minWeight,
				MaxWeight: maxWeight,
			})
			if err != nil {
				return nil, errors.New("get adjacent error")
			}
		}
	}
	return nil, nil
}

func (c CacheGraphRepository) loadVertex(ctx context.Context, key string) (model.Vertex, bool) {
	value, found := c.client.Get(key)
	if !found {
		return nil, found
	}
	switch v := value.(type) {
	case model.Vertex:
		return v, true
	default:
		return nil, false
	}

}

func (c CacheGraphRepository) DumpVertex(ctx context.Context, vertex model.Vertex) error {
	if err := c.client.Add(vertex.Digest(), vertex, c.ttl); err != nil {
		return err
	}
	return nil
}

func (c CacheGraphRepository) DumpEdge(ctx context.Context, edge model.Edge) error {
	if err := c.DumpVertex(ctx, edge.Tail); err != nil {
		return err
	}
	if err := c.DumpVertex(ctx, edge.Head); err != nil {
		return err
	}
	if err := c.client.Add(getKeyOfEdge(edge), edge.Weight, c.ttl); err != nil {
		return err
	}
	return nil
}

func splitToHeadAndTail(edgeKey string) (string, string) {
	keys := strings.Split(edgeKey, "->")
	return keys[0], keys[1]
}

func getKeyOfEdge(edge model.Edge) string {
	return edge.Tail.Digest() + "->" + edge.Head.Digest()
}
