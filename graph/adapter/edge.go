package adapter

import (
	. "github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func LanternEdge(protoEdge *pb.Edge) Edge {
	return Edge{
		Tail:       Key(protoEdge.Tail),
		Head:       Key(protoEdge.Head),
		Weight:     Weight(protoEdge.Weight),
		Expiration: Expiration(protoEdge.Expiration.AsTime().Unix()),
	}
}

func ProtoEdge(LanternEdge Edge) *pb.Edge {
	return &pb.Edge{
		Tail:       string(LanternEdge.Tail),
		Head:       string(LanternEdge.Head),
		Weight:     float32(LanternEdge.Weight),
		Expiration: timestamppb.New(time.Unix(int64(LanternEdge.Expiration), 0)),
	}
}
