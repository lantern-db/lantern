package adapter

import (
	. "github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/pb"
	"time"
)

type ProtoEdge struct {
	message *pb.Edge
}

func (p ProtoEdge) Tail() Key {
	return Key(p.message.Tail)
}

func (p ProtoEdge) Head() Key {
	return Key(p.message.Head)
}

func (p ProtoEdge) Weight() Weight {
	return Weight(p.message.Weight)
}

func (p ProtoEdge) Expiration() Expiration {
	return Expiration(p.message.Expiration.AsTime().Unix())
}

func (p ProtoEdge) AsProto() *pb.Edge {
	return p.message
}

func NewProtoEdge(message *pb.Edge) Edge {
	return ProtoEdge{message: message}
}

func NewProtoEdgeOf(tail Key, head Key, weight Weight, ttl time.Duration) Edge {
	message := &pb.Edge{
		Tail:       string(tail),
		Head:       string(head),
		Weight:     float32(weight),
		Expiration: NewExpiration(ttl).AsProtoTimestamp(),
	}

	return NewProtoEdge(message)
}
