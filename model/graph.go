package model

type Vertex interface {
	Digest() string
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
