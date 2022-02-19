package adapter

import (
	"github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/pb"
)

func LanternGraph(protoGraph *pb.Graph) model.Graph {
	vertexMap := make(model.VertexMap)
	edgeMap := make(model.EdgeMap)
	dfMap := make(model.DocumentFrequency)

	for key, protoVertex := range protoGraph.VertexMap {
		vertexMap[model.Key(key)] = ProtoVertex{protoVertex}
		dfMap[model.Key(key)] = protoGraph.DfMap[key]
	}

	for tail, neighbor := range protoGraph.NeighborMap {
		edgeMap[model.Key(tail)] = make(map[model.Key]model.Edge)
		for head, protoEdge := range neighbor.EdgeMap {
			edgeMap[model.Key(tail)][model.Key(head)] = ProtoEdge{protoEdge}
		}
	}
	return model.Graph{
		VertexMap: vertexMap,
		EdgeMap:   edgeMap,
		Df:        dfMap,
	}
}

func ProtoGraph(graph model.Graph) *pb.Graph {
	g := pb.Graph{}
	g.VertexMap = make(map[string]*pb.Vertex)
	g.DfMap = make(map[string]uint32)
	for key, vertex := range graph.VertexMap {
		g.VertexMap[string(key)] = vertex.AsProto()
		g.DfMap[string(key)] = graph.Df[key]
	}

	g.NeighborMap = make(map[string]*pb.Neighbor)
	for tailKey, heads := range graph.EdgeMap {
		neighbor := pb.Neighbor{}
		neighbor.EdgeMap = make(map[string]*pb.Edge)
		for headKey, edge := range heads {
			neighbor.EdgeMap[string(headKey)] = edge.AsProto()
		}
		g.NeighborMap[string(tailKey)] = &neighbor
	}

	return &g
}
