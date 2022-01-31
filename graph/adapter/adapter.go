package adapter

import (
	. "github.com/lantern-db/lantern/graph/model"
	pb "github.com/lantern-db/lantern/pb"
)

func LanternQuery(request *pb.IlluminateRequest) LoadQuery {
	return LoadQuery{
		Seed:      Key(request.Seed),
		Step:      request.Step,
		MinWeight: request.MinWeight,
		MaxWeight: request.MaxWeight,
	}
}

func LanternGraph(protoGraph *pb.Graph) Graph {
	vertexMap := make(VertexMap)
	edgeMap := make(EdgeMap)

	for key, protoVertex := range protoGraph.VertexMap {
		vertexMap[Key(key)] = LanternVertex(protoVertex)
	}

	for tail, neighbor := range protoGraph.NeighborMap {
		edgeMap[Key(tail)] = make(map[Key]Edge)
		for head, protoEdge := range neighbor.EdgeMap {
			edgeMap[Key(tail)][Key(head)] = LanternEdge(protoEdge)
		}
	}
	return Graph{
		VertexMap: vertexMap,
		EdgeMap:   edgeMap,
	}
}

func ProtoGraph(graph Graph) *pb.Graph {
	g := pb.Graph{}
	g.VertexMap = make(map[string]*pb.Vertex)
	for _, vertex := range graph.VertexMap {
		switch v := vertex.Value.(type) {
		case *pb.Vertex:
			g.VertexMap[string(vertex.Key)] = v

		default:
			g.VertexMap[string(vertex.Key)] = ProtoVertex(vertex)
		}
	}

	g.NeighborMap = make(map[string]*pb.Neighbor)
	for tailKey, heads := range graph.EdgeMap {
		neighbor := pb.Neighbor{}
		neighbor.EdgeMap = make(map[string]*pb.Edge)
		for headKey, edge := range heads {
			neighbor.EdgeMap[string(headKey)] = ProtoEdge(edge)
		}
		g.NeighborMap[string(tailKey)] = &neighbor
	}
	return &g
}
