package service

import (
	"context"
	"github.com/piroyoung/lanterne/adapter"
	"github.com/piroyoung/lanterne/graph/cache"
	pb "github.com/piroyoung/lanterne/grpc"
)

type LanterneService struct {
	cache *cache.GraphCache
}

func (l *LanterneService) Illuminate(ctx context.Context, request *pb.IlluminateRequest) (*pb.IlluminateResponse, error) {
	q := adapter.LanterneQuery(request)
	graph := l.cache.Load(q)
	response := pb.IlluminateResponse{
		Graph: adapter.ProtoGraph(graph),
	}

	return &response, nil
}

func (l LanterneService) DumpVertex(ctx context.Context, vertex *pb.Vertex) (*pb.DumpResponse, error) {
	l.cache.DumpVertex(adapter.LanterneVertex(vertex))
	return &pb.DumpResponse{}, nil
}

func (l LanterneService) DumpEdge(ctx context.Context, edge *pb.Edge) (*pb.DumpResponse, error) {
	l.cache.DumpEdge(adapter.LanterneEdge(edge))
	return &pb.DumpResponse{}, nil
}

func (l LanterneService) mustEmbedUnimplementedLanterneServer() {
}
