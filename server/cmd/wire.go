//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/lantern-db/lantern/graph/cache"
	"github.com/lantern-db/lantern/graph/model"
	"github.com/lantern-db/lantern/server/service"
	"google.golang.org/grpc"
	"net"
)

func newVertexCache(config *model.LanternServerConfig) *cache.VertexCache {
	return cache.NewVertexCache(config.Ttl)
}

func newEdgeCache(config *model.LanternServerConfig) *cache.EdgeCache {
	return cache.NewEdgeCache(config.Ttl)
}

func newGraphCache(v *cache.VertexCache, e *cache.EdgeCache) *cache.GraphCache {
	return cache.NewGraphCache(v, e)
}

func newListener(config *model.LanternServerConfig) (net.Listener, error) {
	return net.Listen("tcp", ":"+config.Port)
}

func newGrpcServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{}
}

func newLanternServer(config *model.LanternServerConfig, listener net.Listener, svc *service.LanternService, server *grpc.Server) *service.LanternServer {
	return service.NewLanternServer(config.FlushInterval, listener, svc, server)
}

func initializeLanternServer() (*service.LanternServer, error) {
	wire.Build(
		model.LoadServerConfig,
		newVertexCache,
		newEdgeCache,
		newGraphCache,
		service.NewLanternService,
		newListener,
		newGrpcServerOptions,
		grpc.NewServer,
		newLanternServer,
	)

	return &service.LanternServer{}, nil
}
