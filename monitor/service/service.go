package service

import (
	"context"
	"errors"
	"github.com/lantern-db/lantern/monitor/config"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusService struct {
	port string
}

func NewPrometheusService(config *config.PrometheusConfig) *PrometheusService {
	return &PrometheusService{config.Port}
}

func (s *PrometheusService) Port() string {
	return ":" + s.port
}

func (s *PrometheusService) Run(ctx context.Context) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{Addr: s.Port(), Handler: mux}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); errors.Is(err, context.Canceled) {
			log.Println("stop metrics server gracefully")
		} else {
			log.Fatal(err)
		}
	}()
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
}
