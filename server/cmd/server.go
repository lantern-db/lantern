package main

import (
	"github.com/piroyoung/lanterne/graph/cache"
	pb "github.com/piroyoung/lanterne/grpc"
	"github.com/piroyoung/lanterne/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	lanternePort := os.Getenv("LANTERNE_PORT")

	ttl, err := strconv.Atoi(os.Getenv("TTL_SECOND"))
	if err != nil {
		log.Fatalf("ttl parse error: %v", err)
	}

	lis, err := net.Listen("tcp", ":"+lanternePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	graphCache := cache.NewGraphCache(time.Duration(ttl) * time.Second)
	svc := service.NewLanterneService(&graphCache)

	s := grpc.NewServer()
	pb.RegisterLanterneServer(s, svc)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
