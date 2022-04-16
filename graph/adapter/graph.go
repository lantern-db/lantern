package adapter

import (
	pb "github.com/lantern-db/lantern-proto/go/lantern/v1"
	m "github.com/lantern-db/lantern/graph/model"
)

func LanternGraph(protoGraph *pb.Graph) m.Graph {
	vertexMap := make(m.VertexMap)
	edgeMap := make(m.EdgeMap)
	vertexStats := make(map[m.Key]m.VertexStats)

	for _, v := range protoGraph.Vertices {
		key := m.Key(v.Key)

		vertexMap[key] = ProtoVertex{v}
	}

	for _, vStat := range protoGraph.Stats.VertexDegrees {
		key := m.Key(vStat.Key)
		vertexStats[key] = m.VertexStats{
			Degree: m.Degree{
				In:  vStat.In,
				Out: vStat.Out,
			},
		}
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
		Stats: m.GraphStats{
			VertexStats: vertexStats,
		},
	}
}

func ProtoGraph(graph m.Graph) *pb.Graph {
	g := &pb.Graph{
		Vertices: []*pb.Vertex{},
		Edges:    []*pb.Edge{},
		Stats: &pb.GraphStats{
			VertexDegrees: []*pb.GraphStats_VertexDegree{},
		},
	}

	for key, vertex := range graph.VertexMap {
		v := vertex.AsProto()
		g.Stats.VertexDegrees = append(g.Stats.VertexDegrees, &pb.GraphStats_VertexDegree{
			Key: string(key),
			In:  graph.Stats.VertexStats[key].Degree.In,
			Out: graph.Stats.VertexStats[key].Degree.Out,
		})
		g.Vertices = append(g.Vertices, v)
	}

	for _, heads := range graph.EdgeMap {
		for _, edge := range heads {
			g.Edges = append(g.Edges, edge.AsProto())
		}
	}

	return g
}
