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

func (g *Graph) Edges() []Edge {
	var edges []Edge
	for tail, heads := range g.EdgeMap {
		for head, weight := range heads {
			edges = append(edges, Edge{
				Tail:   g.VertexMap[tail],
				Head:   g.VertexMap[head],
				Weight: weight,
			})
		}
	}
	return edges
}

func (g *Graph) Vertices() []Vertex {
	var vertices []Vertex

	for _, v := range g.VertexMap {
		vertices = append(vertices, v)
	}

	return vertices
}