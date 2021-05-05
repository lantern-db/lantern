package client

import (
	"github.com/piroyoung/lanterne/graph/model"
	"github.com/piroyoung/lanterne/pb"
)

type VertexMap map[string]*model.ProtoVertex
type NeighborMap map[string]map[string]float32

type IlluminateResult struct {
	VertexMap   VertexMap   `json:"vertexMap,omitempty"`
	NeighborMap NeighborMap `json:"neighborMap,omitempty"`
}

func NewIlluminateResult(graph *pb.Graph) IlluminateResult {
	vertexMap := make(VertexMap)
	for key, value := range graph.VertexMap {
		vertexMap[key] = &model.ProtoVertex{Message: value}
	}
	neighborMap := make(NeighborMap)
	for tailKey, heads := range graph.NeighborMap {
		neighborMap[tailKey] = make(map[string]float32)
		for headKey, weight := range heads.WeightMap {
			neighborMap[tailKey][headKey] = weight
		}
	}

	return IlluminateResult{
		VertexMap:   vertexMap,
		NeighborMap: neighborMap,
	}
}
