package service

import (
	"context"
	"github.com/lantern-db/lantern/graph/adapter"
	"github.com/lantern-db/lantern/graph/cache"
	m "github.com/lantern-db/lantern/graph/model"
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

func (l *LanternService) GetVertex(ctx context.Context, keys *pb.Keys) (*pb.Graph, error) {
	var vertices []*pb.Vertex
	for _, key := range keys.Keys {
		if vertex, ok := l.cache.LoadVertex(m.Key(key)); ok {
			vertices = append(vertices, vertex.AsProto())
		}
	}
	return &pb.Graph{Vertices: vertices}, nil
}

func (l *LanternService) PutVertex(ctx context.Context, graph *pb.Graph) (*pb.PutResponse, error) {
	for _, vertex := range graph.Vertices {
		l.cache.DumpVertex(adapter.NewProtoVertex(vertex))
	}
	return &pb.PutResponse{}, nil
}

func (l *LanternService) PutEdge(ctx context.Context, graph *pb.Graph) (*pb.PutResponse, error) {
	for _, edge := range graph.Edges {
		l.cache.DumpEdge(adapter.NewProtoEdge(edge))
	}
	return &pb.PutResponse{}, nil
}

func (l *LanternService) Put(ctx context.Context, graph *pb.Graph) (*pb.PutResponse, error) {
	for _, vertex := range graph.Vertices {
		l.cache.DumpVertex(adapter.NewProtoVertex(vertex))
	}
	for _, edge := range graph.Edges {
		l.cache.DumpEdge(adapter.NewProtoEdge(edge))
	}
	return &pb.PutResponse{}, nil
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
