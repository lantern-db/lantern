package adapter

import (
	"github.com/piroyoung/lanterne/graph/model"
	pb "github.com/piroyoung/lanterne/pb"
)

func LanterneQuery(request *pb.IlluminateRequest) model.LoadQuery {
	return model.LoadQuery{
		Seed:      &model.ProtoVertex{Message: request.Seed},
		Step:      request.Step,
		MinWeight: request.MinWeight,
		MaxWeight: request.MaxWeight,
	}
}

func LanterneVertex(vertex *pb.Vertex) model.Vertex {
	return &model.ProtoVertex{
		Message: vertex,
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
	g.VertexMap = make(map[string]*pb.Vertex)
	for _, value := range graph.VertexMap {
		switch v := value.Value().(type) {
		case *pb.Vertex:
			g.VertexMap[value.Key()] = v

		default:
			g.VertexMap[value.Key()] = &pb.Vertex{
				Key:   value.Key(),
				Value: &pb.Vertex_Nil{Nil: true},
			}
		}
	}

	g.NeighborMap = make(map[string]*pb.Neighbor)
	for tailKey, heads := range graph.Adjacency {
		neighbor := pb.Neighbor{}
		neighbor.WeightMap = make(map[string]float32)
		for headKey, weight := range heads {
			neighbor.WeightMap[headKey] = weight
		}
		g.NeighborMap[tailKey] = &neighbor
	}
	return &g
}
