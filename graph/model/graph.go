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
	VertexMap map[string]Vertex
	EdgeMap   map[string]map[string]float32
}

func NewGraph() Graph {
	return Graph{
		VertexMap: make(map[string]Vertex),
		EdgeMap:   make(map[string]map[string]float32),
	}
}
