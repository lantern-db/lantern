package adapter

import (
	. "github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func LanternVertex(protoVertex *pb.Vertex) Vertex {
	lanternVertex := Vertex{
		Key:        Key(protoVertex.Key),
		Expiration: Expiration(protoVertex.Expiration.AsTime().Unix()),
	}
	switch v := protoVertex.Value.(type) {
	case *pb.Vertex_Int32:
		lanternVertex.Value = v.Int32

	case *pb.Vertex_Uint32:
		lanternVertex.Value = v.Uint32

	case *pb.Vertex_Int64:
		lanternVertex.Value = v.Int64

	case *pb.Vertex_Uint64:
		lanternVertex.Value = v.Uint64

	case *pb.Vertex_Float32:
		lanternVertex.Value = v.Float32

	case *pb.Vertex_Float64:
		lanternVertex.Value = v.Float64

	case *pb.Vertex_Bool:
		lanternVertex.Value = v.Bool

	case *pb.Vertex_String_:
		lanternVertex.Value = v.String_

	case *pb.Vertex_Bytes:
		lanternVertex.Value = v.Bytes

	case *pb.Vertex_Timestamp:
		lanternVertex.Value = v.Timestamp

	case *pb.Vertex_Nil:
		lanternVertex.Value = nil

	default:
		lanternVertex.Value = nil
	}

	return lanternVertex
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

	case *pb.Vertex:
		protoVertex = v

	default:
		protoVertex.Value = &pb.Vertex_Nil{Nil: true}
	}
	return protoVertex
}
