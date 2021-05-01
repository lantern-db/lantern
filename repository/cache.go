package repository

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/piroyoung/lanterne/model"
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
	panic("implement me")
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

func getKeyOfEdge(edge model.Edge) string {
	return edge.Tail.Digest() + "->" + edge.Head.Digest()
}
