package adapter

import (
	m "github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/pb"
)

func LanternGraph(protoGraph *pb.Graph) m.Graph {
	vertexMap := make(m.VertexMap)
	edgeMap := make(m.EdgeMap)
	dfMap := make(m.DocumentFrequency)

	for _, v := range protoGraph.Vertices {
		key := m.Key(v.Key)

		vertexMap[key] = ProtoVertex{v}
		dfMap[key] = v.InDegree
	}

	for _, e := range protoGraph.Edges {
		tail := m.Key(e.Tail)
		head := m.Key(e.Head)

		if _, ok := edgeMap[tail]; !ok {
			edgeMap[tail] = make(map[m.Key]m.Edge)
		}

		edgeMap[tail][head] = ProtoEdge{e}
	}
	return m.Graph{
		VertexMap: vertexMap,
		EdgeMap:   edgeMap,
		Df:        dfMap,
	}
}

func ProtoGraph(graph m.Graph) *pb.Graph {
	g := &pb.Graph{
		Vertices: []*pb.Vertex{},
		Edges:    []*pb.Edge{},
	}

	for key, vertex := range graph.VertexMap {
		v := vertex.AsProto()
		v.InDegree = graph.Df[key]
		g.Vertices = append(g.Vertices, v)
	}

	for _, heads := range graph.EdgeMap {
		for _, edge := range heads {
			g.Edges = append(g.Edges, edge.AsProto())
		}
	}

	return g
}
