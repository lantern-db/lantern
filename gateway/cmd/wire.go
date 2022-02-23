//+build wireinject

package main

import (
	"flag"
	"github.com/google/wire"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lantern-db/lantern/gateway/config"
	"github.com/lantern-db/lantern/gateway/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newGrpcServerEndpoint(config *config.GatewayConfig) service.EndpointString {
	return flag.String("grpc-server-endpoint", config.LanternHost+":"+config.LanternPort, "gRPC server endpoint")
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
