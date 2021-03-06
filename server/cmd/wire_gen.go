// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/lantern-db/lantern/graph/cache"
	config2 "github.com/lantern-db/lantern/monitor/config"
	service2 "github.com/lantern-db/lantern/monitor/service"
	"github.com/lantern-db/lantern/server/config"
	"github.com/lantern-db/lantern/server/service"
	"google.golang.org/grpc"
	"net"
)

// Injectors from wire.go:

func initializeLanternServer() (*service.LanternServer, error) {
	lanternServerConfig, err := config.LoadServerConfig()
	if err != nil {
		return nil, err
	}
	listener, err := newListener(lanternServerConfig)
	if err != nil {
		return nil, err
	}
	vertexCache := newVertexCache()
	edgeCache := newEdgeCache()
	graphCache := newGraphCache(vertexCache, edgeCache)
	lanternService := service.NewLanternService(graphCache)
	v := newGrpcServerOptions()
	server := grpc.NewServer(v...)
	lanternServer := newLanternServer(lanternServerConfig, listener, lanternService, server)
	return lanternServer, nil
}

func initializePrometheusService() (*service2.PrometheusService, error) {
	prometheusConfig, err := config2.LoadPrometheusConfig()
	if err != nil {
		return nil, err
	}
	prometheusService := service2.NewPrometheusService(prometheusConfig)
	return prometheusService, nil
}

// wire.go:

func newVertexCache() *cache.VertexCache {
	return cache.NewVertexCache()
}

func newEdgeCache() *cache.EdgeCache {
	return cache.NewEdgeCache()
}

func newGraphCache(v *cache.VertexCache, e *cache.EdgeCache) *cache.GraphCache {
	return cache.NewGraphCache(v, e)
}

func newListener(config3 *config.LanternServerConfig) (net.Listener, error) {
	return net.Listen("tcp", ":"+config3.Port)
}

func newGrpcServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{}
}

func newLanternServer(config3 *config.LanternServerConfig, listener net.Listener, svc *service.LanternService, server *grpc.Server) *service.LanternServer {
	return service.NewLanternServer(config3.FlushInterval, listener, svc, server)
}
