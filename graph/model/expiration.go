package model

import "time"

type Expiration int64

func (e Expiration) Dead() bool {
	return time.Now().Unix() > int64(e)
}

func NewExpiration(ttl time.Duration) Expiration {
	return Expiration(time.Now().Add(ttl).Unix())
}
