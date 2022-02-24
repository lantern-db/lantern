//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/lantern-db/lantern/graph/cache"
	"github.com/lantern-db/lantern/server/config"
	"github.com/lantern-db/lantern/server/service"
	"google.golang.org/grpc"
	"net"
)

func newVertexCache() *cache.VertexCache {
	return cache.NewVertexCache()
}

func newEdgeCache() *cache.EdgeCache {
	return cache.NewEdgeCache()
}

func newGraphCache(v *cache.VertexCache, e *cache.EdgeCache) *cache.GraphCache {
	return cache.NewGraphCache(v, e)
}

func newListener(config *config.LanternServerConfig) (net.Listener, error) {
	return net.Listen("tcp", ":"+config.Port)
}

func newGrpcServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{}
}

func newLanternServer(config *config.LanternServerConfig, listener net.Listener, svc *service.LanternService, server *grpc.Server) *service.LanternServer {
	return service.NewLanternServer(config.FlushInterval, listener, svc, server)
}

func initializeLanternServer() (*service.LanternServer, error) {
	wire.Build(
		config.LoadServerConfig,
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
