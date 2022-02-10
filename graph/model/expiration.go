package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Expiration int64

func NewExpiration(ttl time.Duration) Expiration {
	return Expiration(time.Now().Add(ttl).Unix())
}

func (e Expiration) Dead() bool {
	return time.Now().Unix() > int64(e)
}

func (e Expiration) AsTime() time.Time {
	return time.Unix(int64(e), 0)
}

func (e Expiration) AsProtoTimestamp() *timestamppb.Timestamp {
	return timestamppb.New(e.AsTime())
}
