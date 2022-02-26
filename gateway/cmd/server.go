package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	for _, pair := range os.Environ() {
		log.Println(pair)
	}

	server, err := initializeGrpcGatewayServer()
	if err != nil {
		log.Fatal(err)
	}

	if prom, err := initializePrometheusService(); err != nil {
		log.Fatal(err)
	} else {
		prom.Run(ctx)
	}

	signalCh := make(chan os.Signal)
	defer close(signalCh)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	go func() {
		<-signalCh
		cancel()
	}()
	flag.Parse()
	defer glog.Flush()

	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
