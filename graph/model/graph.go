package model

import (
	"errors"
	"github.com/lantern-db/lantern/pb"
	"time"
)

type Key string
type Value interface{}
type Weight float32

type Vertex struct {
	Key        Key        `json:"key,omitempty"`
	Value      Value      `json:"value,omitempty"`
	Expiration Expiration `json:"value,omitempty""`
}

func (p *Vertex) IntValue() (int, error) {
	switch v := p.Value.(type) {
	case *pb.Vertex_Int32:
		return int(v.Int32), nil

	case *pb.Vertex_Uint32:
		return int(v.Uint32), nil

	default:
		return 0, errors.New("parse error int")

	}
}

func (p *Vertex) Int32Value() (int32, error) {
	switch v := p.Value.(type) {
	case *pb.Vertex_Int32:
		return v.Int32, nil

	case *pb.Vertex_Uint32:
		return int32(v.Uint32), nil

	default:
		return 0, errors.New("parse error int64")

	}
}

func (p *Vertex) Int64Value() (int64, error) {
	switch v := p.Value.(type) {
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

func (p *Vertex) Uint32Value() (uint32, error) {
	switch v := p.Value.(type) {
	case *pb.Vertex_Uint32:
		return v.Uint32, nil

	default:
		return 0, errors.New("parse error uint32")

	}
}

func (p *Vertex) Uint64Value() (uint64, error) {
	switch v := p.Value.(type) {
	case *pb.Vertex_Uint32:
		return uint64(v.Uint32), nil

	case *pb.Vertex_Uint64:
		return v.Uint64, nil

	default:
		return 0, errors.New("parse error uint32")

	}
}

func (p *Vertex) Float32Value() (float32, error) {
	switch v := p.Value.(type) {
	case *pb.Vertex_Float32:
		return v.Float32, nil

	default:
		return 0.0, errors.New("parse error float32")

	}
}

func (p *Vertex) Float64Value() (float64, error) {
	switch v := p.Value.(type) {
	case *pb.Vertex_Float32:
		return float64(v.Float32), nil

	case *pb.Vertex_Float64:
		return v.Float64, nil

	default:
		return 0.0, errors.New("parse error float64")

	}
}

func (p *Vertex) BoolValue() (bool, error) {
	switch v := p.Value.(type) {
	case *pb.Vertex_Bool:
		return v.Bool, nil

	default:
		return false, errors.New("parse error bool")

	}
}

func (p *Vertex) StringValue() (string, error) {
	switch v := p.Value.(type) {
	case *pb.Vertex_String_:
		return v.String_, nil

	default:
		return "", errors.New("parse error string")

	}
}

func (p *Vertex) BytesValue() ([]byte, error) {
	switch v := p.Value.(type) {
	case *pb.Vertex_Bytes:
		return v.Bytes, nil

	default:
		return nil, errors.New("parse error bytes")

	}
}

func (p *Vertex) TimeValue() (time.Time, error) {
	switch v := p.Value.(type) {
	case *pb.Vertex_Timestamp:
		return v.Timestamp.AsTime(), nil

	default:
		return time.Now(), errors.New("parse error timestamp")
	}
}

func (p *Vertex) NilValue() (interface{}, error) {
	switch p.Value.(type) {
	case *pb.Vertex_Nil:
		return nil, nil

	default:
		return nil, errors.New("parse error nil")

	}
}

type Edge struct {
	Tail       Key        `json:"tail,omitempty"`
	Head       Key        `json:"head,omitempty"`
	Weight     Weight     `json:"weight,omitempty"`
	Expiration Expiration `json:"expiration,omitempty"`
}

type VertexMap map[Key]Vertex
type EdgeMap map[Key]map[Key]Edge

type Graph struct {
	VertexMap VertexMap `json:"vertexMap,omitempty"`
	EdgeMap   EdgeMap   `json:"edgeMap,omitempty"`
}

func NewGraph() Graph {
	return Graph{
		VertexMap: make(VertexMap),
		EdgeMap:   make(EdgeMap),
	}
}

func (g *Graph) Edges() []Edge {
	var edges []Edge
	for _, heads := range g.EdgeMap {
		for _, edge := range heads {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (g *Graph) Vertices() []Vertex {
	var vertices []Vertex

	for _, v := range g.VertexMap {
		vertices = append(vertices, v)
	}

	return vertices
}
