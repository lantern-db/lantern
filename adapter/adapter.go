package adapter

import (
	"github.com/piroyoung/lanterne/graph/model"
	"github.com/piroyoung/lanterne/grpc"
)

type ProtoVertex struct {
	key string
}

func (j ProtoVertex) Digest() string {
	return j.key
}

func LanternQuery(request *grpc.IlluminateRequest) model.LoadQuery {
	return model.LoadQuery{
		Seed:      ProtoVertex{key: request.Seed.Key},
		Degree:    request.Degree,
		MinWeight: request.MinWeight,
		MaxWeight: request.MaxWeight,
	}
}

func GrpcGraph(graph model.Graph) *grpc.Graph {
	g := grpc.Graph{}
	for _, v := range graph.VertexMap {
		g.Vertices = append(g.Vertices, &grpc.Vertex{
			Key: v.Digest(),
		})
	}

	for tailKey, heads := range graph.EdgeMap {
		for headKey, weight := range heads {
			g.Edges = append(g.Edges, &grpc.Edge{
				Tail:   &grpc.Vertex{Key: tailKey},
				Head:   &grpc.Vertex{Key: headKey},
				Weight: weight,
			})
		}
	}
	return &g
}
