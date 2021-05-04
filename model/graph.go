package model

type Vertex interface {
	Key() string
}

type Edge struct {
	Tail   *Vertex
	Head   *Vertex
	Weight float32
}

type Graph struct {
	Vertices []*Vertex
	Edges    []Edge
}
