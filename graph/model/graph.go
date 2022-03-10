package model

type Key string
type Value interface{}
type Weight float32

type VertexMap map[Key]Vertex
type EdgeMap map[Key]map[Key]Edge

type Degree struct {
	In  uint32 `json:"in,omitempty"`
	Out uint32 `json:"out,omitempty"`
}

type VertexStats struct {
	Degree Degree `json:"degree,omitempty"`
}

type GraphStats struct {
	VertexStats map[Key]VertexStats `json:"vertex,omitempty"`
}

type Graph struct {
	VertexMap VertexMap
	EdgeMap   EdgeMap
	Stats     GraphStats
}

func NewGraph() Graph {
	return Graph{
		VertexMap: make(VertexMap),
		EdgeMap:   make(EdgeMap),
		Stats: GraphStats{
			VertexStats: make(map[Key]VertexStats),
		},
	}
}

func (g *Graph) Edges() []Edge {
	var edges []Edge
	for _, heads := range g.EdgeMap {
		for _, edge := range heads {
			edges = append(edges, edge)
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

func (g Graph) Render() RenderedGraph {
	return NewRenderedGraphOf(g)
}
