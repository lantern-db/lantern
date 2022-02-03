package model

type Edge interface {
	Tail() Key
	Head() Key
	Weight() Weight
	Expiration() Expiration
}

type EdgeExpression struct {
	Tail       Key        `json:"tail,omitempty"`
	Head       Key        `json:"head,omitempty"`
	Weight     Weight     `json:"weight,omitempty"`
	Expiration Expiration `json:"expiration,omitempty"`
}
