package service

import (
	"context"
	"github.com/lantern-db/lantern/graph/adapter"
	"github.com/lantern-db/lantern/graph/cache"
	"github.com/lantern-db/lantern/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type LanternService struct {
	pb.UnimplementedLanternServer
	cache *cache.GraphCache
}

func NewLanternService(graphCache *cache.GraphCache) *LanternService {
	return &LanternService{cache: graphCache}
}

func (l *LanternService) Illuminate(ctx context.Context, request *pb.IlluminateRequest) (*pb.IlluminateResponse, error) {
	q := adapter.LanternQuery(request)
	graph := l.cache.Load(q)
	response := pb.IlluminateResponse{
		Graph: adapter.ProtoGraph(graph),
	}

	return &response, nil
}

func (l *LanternService) DumpVertex(ctx context.Context, vertex *pb.Vertex) (*pb.DumpResponse, error) {
	l.cache.DumpVertex(adapter.LanternVertex(vertex))
	return &pb.DumpResponse{}, nil
}

func (l *LanternService) DumpEdge(ctx context.Context, edge *pb.Edge) (*pb.DumpResponse, error) {
	l.cache.DumpEdge(adapter.LanternEdge(edge))
	return &pb.DumpResponse{}, nil
}

type LanternServer struct {
	flushInterval time.Duration
	listener      net.Listener
	svc           *LanternService
	server        *grpc.Server
}

func NewLanternServer(flushInterval time.Duration, listener net.Listener, svc *LanternService, server *grpc.Server) *LanternServer {
	return &LanternServer{
		flushInterval: flushInterval,
		listener:      listener,
		svc:           svc,
		server:        server,
	}
}

func (s *LanternServer) Run(ctx context.Context) error {
	pb.RegisterLanternServer(s.server, s.svc)
	go func() {
		t := time.NewTicker(s.flushInterval)
	L:
		for {
			select {
			case <-ctx.Done():
				break L
			case <-t.C:
				log.Println("flush expired cache")
				s.svc.cache.Flush()
			}
		}
	}()
	go func() {
		<-ctx.Done()
		log.Println("stop grpc server gracefully")
		s.server.GracefulStop()
	}()
	return s.server.Serve(s.listener)
}
