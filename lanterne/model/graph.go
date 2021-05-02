package model

type Vertex interface {
	Digest() string
}

type Edge struct {
	Tail   Vertex
	Head   Vertex
	Weight float32
}

type Graph struct {
	Vertices map[string]Vertex
	Edges    map[string]map[string]float32
}

func NewGraph() Graph {
	return Graph{
		Vertices: make(map[string]Vertex),
		Edges:    make(map[string]map[string]float32),
	}
}
