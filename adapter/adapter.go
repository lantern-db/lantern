package adapter

import (
	"github.com/piroyoung/lanterne/graph/model"
	pb "github.com/piroyoung/lanterne/grpc"
)

type KeyVertex struct {
	key string
}

func (k KeyVertex) Digest() string {
	return k.key
}

func LanterneQuery(request *pb.IlluminateRequest) model.LoadQuery {
	return model.LoadQuery{
		Seed:      KeyVertex{key: request.Seed.Key},
		Degree:    request.Degree,
		MinWeight: request.MinWeight,
		MaxWeight: request.MaxWeight,
	}
}

func LanterneVertex(vertex *pb.Vertex) model.Vertex {
	return &KeyVertex{
		key: vertex.Key,
	}
}

func LanterneEdge(edge *pb.Edge) model.Edge {
	return model.Edge{
		Tail:   LanterneVertex(edge.Tail),
		Head:   LanterneVertex(edge.Head),
		Weight: edge.Weight,
	}
}

func ProtoGraph(graph model.Graph) *pb.Graph {
	g := pb.Graph{}
	for _, v := range graph.VertexMap {
		g.Vertices = append(g.Vertices, &pb.Vertex{
			Key: v.Digest(),
		})
	}

	for tailKey, heads := range graph.EdgeMap {
		for headKey, weight := range heads {
			g.Edges = append(g.Edges, &pb.Edge{
				Tail:   &pb.Vertex{Key: tailKey},
				Head:   &pb.Vertex{Key: headKey},
				Weight: weight,
			})
		}
	}
	return &g
}
