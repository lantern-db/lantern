package model

import (
	"errors"
	"github.com/piroyoung/lanterne/pb"
	"time"
)

type ProtoVertex struct {
	Message *pb.Vertex
}

func (p *ProtoVertex) Key() string {
	return p.Message.Key
}

func (p *ProtoVertex) Value() interface{} {
	return p.Message
}

func (p *ProtoVertex) Int() (int, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_Int32:
		return int(v.Int32), nil

	case *pb.Vertex_Uint32:
		return int(v.Uint32), nil

	default:
		return 0, errors.New("parse error int")

	}
}

func (p *ProtoVertex) Int32() (int32, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_Int32:
		return v.Int32, nil

	case *pb.Vertex_Uint32:
		return int32(v.Uint32), nil

	default:
		return 0, errors.New("parse error int64")

	}
}

func (p *ProtoVertex) Int64() (int64, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_Int64:
		return v.Int64, nil

	case *pb.Vertex_Uint64:
		return int64(v.Uint64), nil

	case *pb.Vertex_Int32:
		return int64(v.Int32), nil

	case *pb.Vertex_Uint32:
		return int64(v.Uint32), nil

	default:
		return 0, errors.New("parse error int64")

	}
}

func (p *ProtoVertex) Uint32() (uint32, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_Uint32:
		return v.Uint32, nil

	default:
		return 0, errors.New("parse error uint32")

	}
}

func (p *ProtoVertex) Uint64() (uint64, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_Uint32:
		return uint64(v.Uint32), nil

	case *pb.Vertex_Uint64:
		return v.Uint64, nil

	default:
		return 0, errors.New("parse error uint32")

	}
}

func (p *ProtoVertex) Float32() (float32, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_Float32:
		return v.Float32, nil

	default:
		return 0.0, errors.New("parse error float32")

	}
}

func (p *ProtoVertex) Float64() (float64, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_Float32:
		return float64(v.Float32), nil

	case *pb.Vertex_Float64:
		return v.Float64, nil

	default:
		return 0.0, errors.New("parse error float64")

	}
}

func (p *ProtoVertex) Bool() (bool, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_Bool:
		return v.Bool, nil

	default:
		return false, errors.New("parse error bool")

	}
}

func (p *ProtoVertex) String() (string, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_String_:
		return v.String_, nil

	default:
		return "", errors.New("parse error string")

	}
}

func (p *ProtoVertex) Bytes() ([]byte, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_Bytes:
		return v.Bytes, nil

	default:
		return nil, errors.New("parse error bytes")

	}
}

func (p *ProtoVertex) Timestamp() (time.Time, error) {
	switch v := p.Message.Value.(type) {
	case *pb.Vertex_Timestamp:
		return v.Timestamp.AsTime(), nil

	default:
		return time.Now(), errors.New("parse error timestamp")
	}
}

func (p *ProtoVertex) Nil() (interface{}, error) {
	switch p.Message.Value.(type) {
	case *pb.Vertex_Nil:
		return nil, nil

	default:
		return nil, errors.New("parse error nil")

	}
}
