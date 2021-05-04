package model

import "math"

type LoadQuery struct {
	Seed      Vertex
	Step      uint32
	MinWeight float32
	MaxWeight float32
}

func NeighborQuery(seed Vertex, step uint32) LoadQuery {
	return LoadQuery{
		Seed:      seed,
		Step:      step,
		MinWeight: -math.MaxFloat32,
		MaxWeight: math.MaxFloat32,
	}
}

func AdjacentQuery(seed Vertex) LoadQuery {
	return NeighborQuery(seed, 1)
}
