package model

type Edge struct {
	Tail       Key        `json:"tail,omitempty"`
	Head       Key        `json:"head,omitempty"`
	Weight     Weight     `json:"weight,omitempty"`
	Expiration Expiration `json:"expiration,omitempty"`
}
