package model

import "math"

type LoadQuery struct {
	Seed Key
	Step uint32
	TopK uint32
}

func NeighborQuery(seed Key, step uint32, topK uint32) LoadQuery {
	return LoadQuery{
		Seed: seed,
		Step: step,
		TopK: topK,
	}
}

func AdjacentQuery(seed Key) LoadQuery {
	return NeighborQuery(seed, 1, math.MaxUint32)
}
