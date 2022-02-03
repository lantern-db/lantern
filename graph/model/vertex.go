package model

import (
	"github.com/lantern-db/lantern/pb"
	"time"
)

type Vertex interface {
	Key() Key
	Value() Value
	Expiration() Expiration
	AsProto() *pb.Vertex
	StringValue() (string, error)
	IntValue() (int, error)
	Int64Value() (int64, error)
	Float32Value() (float32, error)
	Float64Value() (float64, error)
	BoolValue() (bool, error)
	BytesValue() ([]byte, error)
	TimeValue() (time.Time, error)
	NilValue() (interface{}, error)
}

type VertexExpression struct {
	Key        Key        `json:"key,omitempty"`
	Value      Value      `json:"value,omitempty"`
	Expiration Expiration `json:"expiration,omitempty""`
}
