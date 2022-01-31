package adapter

import (
	. "github.com/lantern-db/lantern/graph/model"
	pb "github.com/lantern-db/lantern/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func LanternQuery(request *pb.IlluminateRequest) LoadQuery {
	return LoadQuery{
		Seed:      Key(request.Seed),
		Step:      request.Step,
		MinWeight: request.MinWeight,
		MaxWeight: request.MaxWeight,
	}
}

func LanternEdge(edge *pb.Edge) Edge {
	return Edge{
		Tail:   Key(edge.Tail),
		Head:   Key(edge.Head),
		Weight: Weight(edge.Weight),
		Expiration: Expiration(edge.Expiration.AsTime().Unix()),
	}
}

func ProtoGraph(graph Graph) *pb.Graph {
	g := pb.Graph{}
	g.VertexMap = make(map[string]*pb.Vertex)
	for _, vertex := range graph.VertexMap {
		switch v := vertex.Value.(type) {
		case *pb.Vertex:
			g.VertexMap[string(vertex.Key)] = v

		default:
			g.VertexMap[string(vertex.Key)] = &pb.Vertex{
				Key:        string(vertex.Key),
				Value:      &pb.Vertex_Nil{Nil: true},
				Expiration: timestamppb.New(time.Unix(int64(vertex.Expiration), 0)),
			}
		}
	}

	g.NeighborMap = make(map[string]*pb.Neighbor)
	for tailKey, heads := range graph.EdgeMap {
		neighbor := pb.Neighbor{}
		neighbor.WeightMap = make(map[string]float32)
		for headKey, weight := range heads {
			neighbor.WeightMap[string(headKey)] = float32(weight)
		}
		g.NeighborMap[string(tailKey)] = &neighbor
	}
	return &g
}
