package service

import (
	"context"
	"github.com/lantern-db/lantern/graph/adapter"
	"github.com/lantern-db/lantern/graph/cache"
	pb "github.com/lantern-db/lantern/pb"
)

type LanternService struct {
	pb.UnimplementedLanternServer
	cache *cache.GraphCache
}

func NewLanternService(graphCache *cache.GraphCache) *LanternService {
	return &LanternService{cache: graphCache}
}

func (l *LanternService) Illuminate(ctx context.Context, request *pb.IlluminateRequest) (*pb.IlluminateResponse, error) {
	q := adapter.LanternQuery(request)
	graph := l.cache.Load(q)
	response := pb.IlluminateResponse{
		Graph: adapter.ProtoGraph(graph),
	}

	return &response, nil
}

func (l *LanternService) DumpVertex(ctx context.Context, vertex *pb.Vertex) (*pb.DumpResponse, error) {
	l.cache.DumpVertex(adapter.LanternVertex(vertex))
	return &pb.DumpResponse{}, nil
}

func (l *LanternService) DumpEdge(ctx context.Context, edge *pb.Edge) (*pb.DumpResponse, error) {
	le := adapter.LanternEdge(edge)

	switch edge.Tail.Value.(type) {
	case *pb.Vertex_Nil:
		v, found := l.cache.LoadVertex(edge.Tail.Key)
		if found {
			le.Tail = v
		}
	}

	switch edge.Head.Value.(type) {
	case *pb.Vertex_Nil:
		v, found := l.cache.LoadVertex(edge.Head.Key)
		if found {
			le.Head = v
		}
	}

	l.cache.DumpEdge(le)
	return &pb.DumpResponse{}, nil
}
