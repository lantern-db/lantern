package model

type NeighborQuery struct {
	Seed      *Vertex
	Step      int
	MinWeight float32
	MaxWeight float32
}

type AdjacentQuery struct {
	Seed *Vertex
}
