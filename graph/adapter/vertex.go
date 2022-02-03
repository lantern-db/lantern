package adapter

import (
	"errors"
	. "github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/pb"
	"time"
)

type ProtoVertex struct {
	message pb.Vertex
}

func (p *ProtoVertex) Key() Key {
	return Key(p.message.Key)
}

func (p *ProtoVertex) Value() Value {
	return p.message
}

func (p *ProtoVertex) AsProto() pb.Vertex {
	return p.message
}

func (p *ProtoVertex) StringValue() (string, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_String_:
		return v.String_, nil
	default:
		return "", errors.New("parse error")
	}
}

func (p *ProtoVertex) Expiration() Expiration {
	return Expiration(p.message.Expiration.AsTime().Unix())
}

func (p *ProtoVertex) IntValue() (int, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Int32:
		return int(v.Int32), nil
	default:
		return 0, errors.New("parse error")
	}
}

func (p *ProtoVertex) Int64Value() (int64, error) {
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

func (p *ProtoVertex) Float32Value() (float32, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Float32:
		return v.Float32, nil
	default:
		return 0.0, errors.New("parse error")
	}
}

func (p *ProtoVertex) Float64Value() (float64, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Float64:
		return v.Float64, nil
	default:
		return 0.0, errors.New("parse error")
	}
}

func (p *ProtoVertex) BoolValue() (bool, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Bool:
		return v.Bool, nil
	default:
		return false, errors.New("parse error")
	}
}

func (p *ProtoVertex) BytesValue() ([]byte, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Bytes:
		return v.Bytes, nil
	default:
		return nil, errors.New("parse error")
	}
}

func (p *ProtoVertex) TimeValue() (time.Time, error) {
	switch v := p.message.Value.(type) {
	case *pb.Vertex_Timestamp:
		return v.Timestamp.AsTime(), nil
	default:
		return time.Unix(0, 0), errors.New("parse error")
	}
}

func (p *ProtoVertex) NilValue() (interface{}, error) {
	return nil, nil
}
