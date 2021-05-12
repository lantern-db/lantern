package main

import (
	"github.com/lantern-db/lantern/graph/cache"
	pb "github.com/lantern-db/lantern/pb"
	"github.com/lantern-db/lantern/server/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	for _, pair := range os.Environ() {
		log.Println(pair)
	}

	flushInterval, err := strconv.Atoi(os.Getenv("LANTERN_FLUSH_INTERVAL"))
	if err != nil {
		log.Fatalf("flush interval parse failed: %v", err)
	}
	LanternPort := os.Getenv("LANTERN_PORT")
	ttl, err := strconv.Atoi(os.Getenv("LANTERN_TTL"))
	if err != nil {
		log.Fatalf("ttl parse error: %v", err)
	}

	signalCh := make(chan os.Signal)
	stopCh := make(chan bool)
	defer close(signalCh)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	graphCache := cache.NewGraphCache(time.Duration(ttl) * time.Second)
	go func() {
		t := time.NewTicker(time.Duration(flushInterval) * time.Second)
	L:
		for {
			select {
			case sig := <-signalCh:
				log.Printf("exit with %v", sig)
				break L
			case <-t.C:
				graphCache.Flush()
			}
		}
		stopCh <- true
	}()

	lis, err := net.Listen("tcp", ":"+LanternPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	svc := service.NewLanternService(&graphCache)
	s := grpc.NewServer()
	pb.RegisterLanternServer(s, svc)

	go func() {
		<-stopCh
		log.Println("stop grpc server gracefully")
		s.GracefulStop()
	}()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
