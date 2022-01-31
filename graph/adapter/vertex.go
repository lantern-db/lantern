package adapter

import (
	. "github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func LanternVertex(protoVertex *pb.Vertex) Vertex {
	return Vertex{
		Key:        Key(protoVertex.Key),
		Expiration: Expiration(protoVertex.Expiration.AsTime().Unix()),
		Value:      protoVertex,
	}
}

func ProtoVertex(lanternVertex Vertex) *pb.Vertex {
	protoVertex := &pb.Vertex{
		Key:        string(lanternVertex.Key),
		Expiration: timestamppb.New(time.Unix(int64(lanternVertex.Expiration), 0)),
	}
	switch v := lanternVertex.Value.(type) {
	case int:
		protoVertex.Value = &pb.Vertex_Int32{Int32: int32(v)}

	case float64:
		protoVertex.Value = &pb.Vertex_Float64{Float64: v}

	case float32:
		protoVertex.Value = &pb.Vertex_Float32{Float32: v}

	case int32:
		protoVertex.Value = &pb.Vertex_Int32{Int32: v}

	case int64:
		protoVertex.Value = &pb.Vertex_Int64{Int64: v}

	case uint32:
		protoVertex.Value = &pb.Vertex_Uint32{Uint32: v}

	case uint64:
		protoVertex.Value = &pb.Vertex_Uint64{Uint64: v}

	case bool:
		protoVertex.Value = &pb.Vertex_Bool{Bool: v}

	case string:
		protoVertex.Value = &pb.Vertex_String_{String_: v}

	case []byte:
		protoVertex.Value = &pb.Vertex_Bytes{Bytes: v}

	case time.Time:
		protoVertex.Value = &pb.Vertex_Timestamp{Timestamp: timestamppb.New(v)}

	case nil:
		protoVertex.Value = &pb.Vertex_Nil{Nil: true}

	default:
		protoVertex.Value = &pb.Vertex_Nil{Nil: true}
	}
	return protoVertex
}
