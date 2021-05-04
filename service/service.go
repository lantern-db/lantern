package service

import (
	"context"
	"github.com/piroyoung/lanterne/adapter"
	"github.com/piroyoung/lanterne/graph/cache"
	pb "github.com/piroyoung/lanterne/grpc"
)

type LanterneService struct {
	pb.UnimplementedLanterneServer
	cache *cache.GraphCache
}

func NewLanterneService(graphCache *cache.GraphCache) *LanterneService {
	return &LanterneService{cache: graphCache}
}

func (l *LanterneService) Illuminate(ctx context.Context, request *pb.IlluminateRequest) (*pb.IlluminateResponse, error) {
	q := adapter.LanterneQuery(request)
	graph := l.cache.Load(q)
	response := pb.IlluminateResponse{
		Graph: adapter.ProtoGraph(graph),
	}

	return &response, nil
}

func (l *LanterneService) DumpVertex(ctx context.Context, vertex *pb.Vertex) (*pb.DumpResponse, error) {
	l.cache.DumpVertex(adapter.LanterneVertex(vertex))
	return &pb.DumpResponse{}, nil
}

func (l *LanterneService) DumpEdge(ctx context.Context, edge *pb.Edge) (*pb.DumpResponse, error) {
	le := adapter.LanterneEdge(edge)

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
