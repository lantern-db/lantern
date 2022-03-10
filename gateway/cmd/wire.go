//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lantern-db/lantern/gateway/config"
	"github.com/lantern-db/lantern/gateway/service"
	promConfig "github.com/lantern-db/lantern/monitor/config"
	promServer "github.com/lantern-db/lantern/monitor/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newGrpcServerEndpoint(config *config.GatewayConfig) service.EndpointString {
	endpointString := config.LanternHost + ":" + config.LanternPort
	return service.EndpointString(&endpointString)
}

func newServeMux() *runtime.ServeMux {
	return runtime.NewServeMux()
}

func newDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

func initializeGrpcGatewayServer() (*service.GrpcGatewayServer, error) {
	wire.Build(
		config.LoadGatewayConfig,
		newGrpcServerEndpoint,
		newServeMux,
		newDialOptions,
		service.NewGrpcGatewayServer,
	)

	return &service.GrpcGatewayServer{}, nil
}

func initializePrometheusService() (*promServer.PrometheusService, error) {
	wire.Build(
		promConfig.LoadPrometheusConfig,
		promServer.NewPrometheusService,
	)

	return &promServer.PrometheusService{}, nil

}
