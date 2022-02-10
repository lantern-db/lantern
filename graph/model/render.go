package model

type RenderedGraph struct {
	Vertices map[Key]Value          `json:"vertices,omitempty"`
	Edges    map[Key]map[Key]Weight `json:"edges,omitempty"`
}

func NewRenderedGraphOf(graph Graph) RenderedGraph {
	vertices := make(map[Key]Value)
	edges := make(map[Key]map[Key]Weight)

	for key, vertex := range graph.VertexMap {
		vertices[key] = vertex.Value()
	}

	for tail, heads := range graph.EdgeMap {
		edges[tail] = make(map[Key]Weight)

		for head, edge := range heads {
			edges[tail][head] = edge.Weight()
		}
	}

	return RenderedGraph{Vertices: vertices, Edges: edges}

}
