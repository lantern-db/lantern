package model

import "math"

type LoadQuery struct {
	Seed      Key
	Step      uint32
	MinWeight float32
	MaxWeight float32
}

func NeighborQuery(seed Key, step uint32) LoadQuery {
	return LoadQuery{
		Seed:      seed,
		Step:      step,
		MinWeight: -math.MaxFloat32,
		MaxWeight: math.MaxFloat32,
	}
}

func AdjacentQuery(seed Key) LoadQuery {
	return NeighborQuery(seed, 1)
}
