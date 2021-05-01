package repository

import (
	"context"
	"github.com/piroyoung/lanterne/model"
)

type InMemoryGraphRepository struct {
	vertexMap map[string]model.Vertex
	edgeMap   map[string]map[string]float32
}

func (i InMemoryGraphRepository) LoadNeighbor(ctx context.Context, query model.NeighborQuery) (model.Graph, error) {
	panic("implement me")
}

func (i InMemoryGraphRepository) LoadAdjacent(ctx context.Context, query model.AdjacentQuery) (model.Graph, error) {
	g := model.NewGraph()
	g.Vertices = append(g.Vertices, query.Seed)

	for headDigest, weight := range i.edgeMap[query.Seed.Digest()] {
		if query.MinWeight <= weight && weight <= query.MaxWeight {
			g.Vertices = append(g.Vertices, i.vertexMap[headDigest])
			g.Edges = append(g.Edges, model.Edge{
				Tail:   query.Seed,
				Head:   i.vertexMap[headDigest],
				Weight: weight,
			})
		}
	}
	return g, nil
}

func (i InMemoryGraphRepository) DumpVertex(ctx context.Context, vertex model.Vertex) error {
	i.vertexMap[vertex.Digest()] = vertex
	return nil
}

func (i InMemoryGraphRepository) DumpEdge(ctx context.Context, edge model.Edge) error {
	if err := i.DumpVertex(ctx, edge.Head); err != nil {
		return err
	}

	if err := i.DumpVertex(ctx, edge.Tail); err != nil {
		return err
	}

	i.edgeMap[edge.Tail.Digest()][edge.Head.Digest()] = edge.Weight
	return nil
}
