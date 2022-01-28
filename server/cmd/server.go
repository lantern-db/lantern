package main

import (
	"context"
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

	server, err := initializeLanternServer()
	if err != nil {
		log.Fatal(err)
	}

	signalCh := make(chan os.Signal)
	defer close(signalCh)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	go func() {
		<-signalCh
		cancel()
	}()
	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
