package model

import "github.com/lantern-db/lantern/pb"

type Edge interface {
	Tail() Key
	Head() Key
	Weight() Weight
	Expiration() Expiration
	AsProto() *pb.Edge
}

type EdgeExpression struct {
	Tail       Key        `json:"tail,omitempty"`
	Head       Key        `json:"head,omitempty"`
	Weight     Weight     `json:"weight,omitempty"`
	Expiration Expiration `json:"expiration,omitempty"`
}
