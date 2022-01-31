package model

import (
	"errors"
	"reflect"
	"time"
)

type Key string
type Value interface{}
type Weight float32

type Vertex struct {
	Key        Key        `json:"key,omitempty"`
	Value      Value      `json:"value,omitempty"`
	Expiration Expiration `json:"expiration,omitempty""`
}

func (v *Vertex) IntValue() (int, error) {
	switch v := v.Value.(type) {
	case int:
		return v, nil

	case int32:
		return int(v), nil

	case uint32:
		return int(v), nil

	default:
		return 0, errors.New("parse error")
	}
}
func (v *Vertex) Int64Value() (int64, error) {
	switch v := v.Value.(type) {
	case int:
		return int64(v), nil

	case int32:
		return int64(v), nil

	case uint32:
		return int64(v), nil

	case int64:
		return v, nil

	case uint64:
		return int64(v), nil

	default:
		return 0, errors.New("parse error")
	}
}

func (v *Vertex) Float32Value() (float32, error) {
	switch v := v.Value.(type) {
	case float32:
		return v, nil

	default:
		return 0.0, errors.New("parse error")
	}
}

func (v *Vertex) Float64Value() (float64, error) {
	switch v := v.Value.(type) {
	case float32:
		return float64(v), nil

	case float64:
		return v, nil

	default:
		return 0.0, errors.New("parse error")
	}
}

func (v *Vertex) BoolValue() (bool, error) {
	switch v := v.Value.(type) {
	case bool:
		return v, nil

	default:
		return false, errors.New("parse error")
	}
}

func (v *Vertex) StringValue() (string, error) {
	switch v := v.Value.(type) {
	case string:
		return v, nil

	default:
		return "", errors.New("parse error")
	}
}

func (v *Vertex) BytesValue() ([]byte, error) {
	switch v := v.Value.(type) {
	case []byte:
		return v, nil

	default:
		return nil, errors.New("parse error")
	}
}

func (v *Vertex) TimeValue() (time.Time, error) {
	switch v := v.Value.(type) {
	case time.Time:
		return v, nil

	default:
		return time.Now(), errors.New("parse error")
	}
}

func (v *Vertex) NilValue() (interface{}, error) {
	if v.Value == nil || reflect.ValueOf(v.Value).IsNil() {
		return nil, nil
	} else {
		return nil, errors.New("parse error")
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
