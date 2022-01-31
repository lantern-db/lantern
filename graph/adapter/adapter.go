package adapter

import (
	"errors"
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

func LanternVertex(vertex *pb.Vertex) Vertex {
	return Vertex{
		Key:        Key(vertex.Key),
		Expiration: Expiration(vertex.Expiration.AsTime().Unix()),
		Value:      vertex.Value,
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

func NewProtoVertex(key string, value interface{}) (*pb.Vertex, error) {
	vertex := &pb.Vertex{
		Key: key,
	}
	switch v := value.(type) {
	case int:
		vertex.Value = &pb.Vertex_Int32{Int32: int32(v)}

	case float64:
		vertex.Value = &pb.Vertex_Float64{Float64: v}

	case float32:
		vertex.Value = &pb.Vertex_Float32{Float32: v}

	case int32:
		vertex.Value = &pb.Vertex_Int32{Int32: v}

	case int64:
		vertex.Value = &pb.Vertex_Int64{Int64: v}

	case uint32:
		vertex.Value = &pb.Vertex_Uint32{Uint32: v}

	case uint64:
		vertex.Value = &pb.Vertex_Uint64{Uint64: v}

	case bool:
		vertex.Value = &pb.Vertex_Bool{Bool: v}

	case string:
		vertex.Value = &pb.Vertex_String_{String_: v}

	case []byte:
		vertex.Value = &pb.Vertex_Bytes{Bytes: v}

	case time.Time:
		vertex.Value = &pb.Vertex_Timestamp{Timestamp: timestamppb.New(v)}

	case nil:
		vertex.Value = &pb.Vertex_Nil{Nil: true}

	default:
		return nil, errors.New("type mismatch")
	}
	return vertex, nil
}
