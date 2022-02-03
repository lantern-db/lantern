package adapter

import (
	. "github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/pb"
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
