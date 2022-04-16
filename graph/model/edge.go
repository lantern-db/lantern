package model

import pb "github.com/lantern-db/lantern-proto/go/lantern/v1"

type Edge interface {
	Tail() Key
	Head() Key
	Weight() Weight
	Expiration() Expiration
	AsProto() *pb.Edge
}

type StaticEdge struct {
	tail       Key
	head       Key
	weight     Weight
	expiration Expiration
}

func NewStaticEdge(tail Key, head Key, weight Weight, expiration Expiration) StaticEdge {
	return StaticEdge{
		tail:       tail,
		head:       head,
		weight:     weight,
		expiration: expiration,
	}
}

func (s StaticEdge) Tail() Key {
	return s.tail
}

func (s StaticEdge) Head() Key {
	return s.head
}

func (s StaticEdge) Weight() Weight {
	return s.weight
}

func (s StaticEdge) Expiration() Expiration {
	return s.expiration
}

func (s StaticEdge) AsProto() *pb.Edge {
	return &pb.Edge{
		Tail:       string(s.tail),
		Head:       string(s.head),
		Weight:     float32(s.weight),
		Expiration: s.expiration.AsProtoTimestamp(),
	}
}
