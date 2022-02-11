package model

import (
	"github.com/lantern-db/lantern/errors"
	"github.com/lantern-db/lantern/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type EmptyVertex struct {
	key        Key
	expiration Expiration
}

func (e EmptyVertex) Key() Key {
	return e.key
}

func (e EmptyVertex) Value() Value {
	return nil
}

func (e EmptyVertex) Expiration() Expiration {
	return e.expiration
}

func (e EmptyVertex) AsProto() *pb.Vertex {
	return &pb.Vertex{
		Key:        string(e.key),
		Expiration: timestamppb.New(e.expiration.AsTime()),
		Value:      &pb.Vertex_Nil{Nil: true},
	}
}

func (e EmptyVertex) StringValue() (string, error) {
	return "", errors.ValueParseError
}

func (e EmptyVertex) IntValue() (int, error) {
	return 0, errors.ValueParseError
}

func (e EmptyVertex) Int64Value() (int64, error) {
	return 0, errors.ValueParseError
}

func (e EmptyVertex) Float32Value() (float32, error) {
	return 0.0, errors.ValueParseError
}

func (e EmptyVertex) Float64Value() (float64, error) {
	return 0.0, errors.ValueParseError
}

func (e EmptyVertex) BoolValue() (bool, error) {
	return false, errors.ValueParseError
}

func (e EmptyVertex) BytesValue() ([]byte, error) {
	return nil, errors.ValueParseError
}

func (e EmptyVertex) TimeValue() (time.Time, error) {
	return time.Now(), errors.ValueParseError
}

func (e EmptyVertex) NilValue() (interface{}, error) {
	return nil, nil
}

func NewEmptyVertexOf(key Key, expiration Expiration) Vertex {
	return EmptyVertex{key: key, expiration: expiration}
}
