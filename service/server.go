package service

import (
	"context"
	"github.com/piroyoung/lanterne/adapter"
	"github.com/piroyoung/lanterne/graph/cache"
	"github.com/piroyoung/lanterne/grpc"
)

type LanterneService struct {
	grpc.UnimplementedLanterneServer
	cache *cache.GraphCache
}

func (s *LanterneService) Illuminate(ctx context.Context, request *grpc.IlluminateRequest) (*grpc.IlluminateResponse, error) {
	q := adapter.LanternQuery(request)
	graph := s.cache.Load(q)
	response := grpc.IlluminateResponse{
		Graph: adapter.GrpcGraph(graph),
	}

	return &response, nil
}
