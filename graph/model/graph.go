package model

type Key string

type Vertex struct {
	Key Key
	Value interface{}
}

type Edge struct {
	Tail   Key
	Head   Key
	Weight float32
}

type Graph struct {
	VertexMap map[Key]Vertex
	EdgeMap   map[Key]map[Key]float32
}

func NewGraph() Graph {
	return Graph{
		VertexMap: make(map[Key]Vertex),
		EdgeMap:   make(map[Key]map[Key]float32),
	}
}

func (g *Graph) Edges() []Edge {
	var edges []Edge
	for tail, heads := range g.EdgeMap {
		for head, weight := range heads {
			edges = append(edges, Edge{
				Tail:   tail,
				Head:   head,
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
