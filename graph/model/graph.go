package model

type Vertex interface {
	Key() string
	Value() interface{}
}

type Edge struct {
	Tail   Vertex
	Head   Vertex
	Weight float32
}

type Graph struct {
	VertexMap map[string]Vertex
	Adjacency map[string]map[string]float32
}

func NewGraph() Graph {
	return Graph{
		VertexMap: make(map[string]Vertex),
		Adjacency: make(map[string]map[string]float32),
	}
}

func (g *Graph) Edges() []Edge {
	var edges []Edge
	for tail, heads := range g.Adjacency {
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
