package service

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/lantern-db/lantern-proto/go/lantern/v1"
	"github.com/lantern-db/lantern/gateway/config"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

type EndpointString *string

type GrpcGatewayServer struct {
	endpoint EndpointString
	mux      *runtime.ServeMux
	opts     []grpc.DialOption
	config   *config.GatewayConfig
}

func NewGrpcGatewayServer(endpoint EndpointString, mux *runtime.ServeMux, opt []grpc.DialOption, config *config.GatewayConfig) *GrpcGatewayServer {
	return &GrpcGatewayServer{
		endpoint: endpoint,
		mux:      mux,
		opts:     opt,
		config:   config,
	}
}

func (s *GrpcGatewayServer) Port() string {
	return ":" + s.config.GatewayPort
}

func (s *GrpcGatewayServer) Run(ctx context.Context) error {
	err := pb.RegisterLanternServiceHandlerFromEndpoint(ctx, s.mux, *s.endpoint, s.opts)
	if err != nil {
		return err
	}

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{s.config.AllowedOrigin},
		AllowedMethods: []string{"GET", "PUT", "DELETE"},
	}).Handler(s.mux)
	srv := &http.Server{Addr: s.Port(), Handler: handler}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); errors.Is(err, context.Canceled) {
			log.Println("stop gateway server gracefully")
		} else {
			log.Fatal(err)
		}
	}()

	if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		return nil
	} else {
		return err
	}
}
