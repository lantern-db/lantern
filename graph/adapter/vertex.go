package adapter

import (
	"errors"
	. "github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/pb"
	"go/types"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ProtoVertex struct {
	message *pb.Vertex
}

func (p ProtoVertex) Key() Key {
	return Key(p.message.Key)
}

func (p ProtoVertex) Value() Value {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Int32:
		return v.Int32

	case *pb.Vertex_Uint32:
		return v.Uint32

	case *pb.Vertex_Int64:
		return v.Int64

	case *pb.Vertex_Uint64:
		return v.Uint64

	case *pb.Vertex_Float32:
		return v.Float32

	case *pb.Vertex_Float64:
		return v.Float64

	case *pb.Vertex_String_:
		return v.String_

	case *pb.Vertex_Bool:
		return v.Bool

	case *pb.Vertex_Timestamp:
		return v.Timestamp.AsTime()

	case *pb.Vertex_Bytes:
		return v.Bytes

	case *pb.Vertex_Nil:
		return nil

	default:
		return v
	}
}

func (p ProtoVertex) AsProto() *pb.Vertex {
	return p.message
}

func (p ProtoVertex) StringValue() (string, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_String_:
		return v.String_, nil
	default:
		return "", errors.New("parse error")
	}
}

func (p ProtoVertex) Expiration() Expiration {
	return Expiration(p.message.Expiration.AsTime().Unix())
}

func (p ProtoVertex) IntValue() (int, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Int32:
		return int(v.Int32), nil
	default:
		return 0, errors.New("parse error")
	}
}

func (p ProtoVertex) Int64Value() (int64, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Int32:
		return int64(v.Int32), nil
	case *pb.Vertex_Uint32:
		return int64(v.Uint32), nil
	case *pb.Vertex_Int64:
		return v.Int64, nil
	case *pb.Vertex_Uint64:
		return int64(v.Uint64), nil
	default:
		return 0, errors.New("parse error")
	}
}

func (p ProtoVertex) Float32Value() (float32, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Float32:
		return v.Float32, nil
	default:
		return 0.0, errors.New("parse error")
	}
}

func (p ProtoVertex) Float64Value() (float64, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Float64:
		return v.Float64, nil
	default:
		return 0.0, errors.New("parse error")
	}
}

func (p ProtoVertex) BoolValue() (bool, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Bool:
		return v.Bool, nil
	default:
		return false, errors.New("parse error")
	}
}

func (p ProtoVertex) BytesValue() ([]byte, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Bytes:
		return v.Bytes, nil
	default:
		return nil, errors.New("parse error")
	}
}

func (p ProtoVertex) TimeValue() (time.Time, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Timestamp:
		return v.Timestamp.AsTime(), nil
	default:
		return time.Unix(0, 0), errors.New("parse error")
	}
}

func (p ProtoVertex) NilValue() (interface{}, error) {
	return nil, nil
}

func NewProtoVertex(message *pb.Vertex) Vertex {
	return ProtoVertex{message: message}
}

func NewProtoVertexOf(key Key, value Value, expiration Expiration) (Vertex, error) {
	message := &pb.Vertex{}
	message.Key = string(key)
	message.Expiration = expiration.AsProtoTimestamp()

	switch v := value.(type) {
	case int:
		message.Value = &pb.Vertex_Int32{Int32: int32(v)}

	case int32:
		message.Value = &pb.Vertex_Int32{Int32: v}

	case int64:
		message.Value = &pb.Vertex_Int64{Int64: v}

	case uint8:
		message.Value = &pb.Vertex_Uint32{Uint32: uint32(v)}

	case uint16:
		message.Value = &pb.Vertex_Uint32{Uint32: uint32(v)}

	case uint32:
		message.Value = &pb.Vertex_Uint32{Uint32: v}

	case uint64:
		message.Value = &pb.Vertex_Uint64{Uint64: v}

	case float32:
		message.Value = &pb.Vertex_Float32{Float32: v}

	case float64:
		message.Value = &pb.Vertex_Float64{Float64: v}

	case bool:
		message.Value = &pb.Vertex_Bool{Bool: v}

	case string:
		message.Value = &pb.Vertex_String_{String_: v}

	case []byte:
		message.Value = &pb.Vertex_Bytes{Bytes: v}

	case time.Time:
		message.Value = &pb.Vertex_Timestamp{Timestamp: timestamppb.New(v)}

	case types.Nil:
		message.Value = &pb.Vertex_Nil{Nil: true}

	default:
		return ProtoVertex{}, errors.New("parse error")

	}

	return NewProtoVertex(message), nil

}
