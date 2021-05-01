package repository

import (
	"context"
	"github.com/piroyoung/lanterne/model"
)

type GraphRepository interface {
	LoadNeighbor(ctx context.Context, query model.NeighborQuery) (model.Graph, error)
	LoadAdjacent(ctx context.Context, query model.AdjacentQuery) (model.Graph, error)
	DumpVertex(ctx context.Context, vertex model.Vertex) error
	DumpEdge(ctx context.Context, edge model.Edge) error
}
