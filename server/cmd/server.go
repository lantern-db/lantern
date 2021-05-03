package main

import (
	"github.com/piroyoung/lanterne/graph/cache"
	pb "github.com/piroyoung/lanterne/grpc"
	"github.com/piroyoung/lanterne/service"
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
	flushInterval, err := strconv.Atoi(os.Getenv("FLUSH_INTERVAL_SECOND"))
	if err != nil {
		log.Fatalf("flush interval parse failed: %v", err)
	}
	lanternePort := os.Getenv("LANTERNE_PORT")
	ttl, err := strconv.Atoi(os.Getenv("TTL_SECOND"))
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

	lis, err := net.Listen("tcp", ":"+lanternePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	svc := service.NewLanterneService(&graphCache)
	s := grpc.NewServer()
	pb.RegisterLanterneServer(s, svc)

	go func() {
		<-stopCh
		log.Println("stop grpc server gracefully")
		s.GracefulStop()
	}()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
