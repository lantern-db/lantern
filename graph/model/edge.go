package model

import "github.com/lantern-db/lantern/pb"

type Edge interface {
	Tail() Key
	Head() Key
	Weight() Weight
	Expiration() Expiration
	AsProto() *pb.Edge
}
