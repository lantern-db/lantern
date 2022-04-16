package service

import (
	"context"
	"errors"
	pb "github.com/lantern-db/lantern-proto/go/graph/v1"
	"github.com/lantern-db/lantern/graph/adapter"
	"github.com/lantern-db/lantern/graph/cache"
	m "github.com/lantern-db/lantern/graph/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
	"time"
)

var (
	optProceed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_request_total",
		Help: "How many grpc request proceeded",
	}, []string{"method", "result"})
)

type LanternService struct {
	pb.UnimplementedLanternServiceServer
	cache *cache.GraphCache
}

func NewLanternService(graphCache *cache.GraphCache) *LanternService {
	return &LanternService{cache: graphCache}
}

func (l *LanternService) Illuminate(ctx context.Context, request *pb.IlluminateRequest) (*pb.IlluminateResponse, error) {
	if request.TopK == 0 {
		request.TopK = math.MaxInt32
	}

	q := adapter.LanternQuery(request)
	graph := l.cache.Load(q)
	response := pb.IlluminateResponse{
		Graph:  adapter.ProtoGraph(graph),
		Status: pb.Status_STATUS_OK,
	}
	optProceed.WithLabelValues("illuminate", "success").Inc()
	return &response, nil
}

func (l *LanternService) GetVertex(ctx context.Context, request *pb.GetVertexRequest) (*pb.GetVertexResponse, error) {
	if vertex, ok := l.cache.GetVertex(m.Key(request.Key)); ok {
		return &pb.GetVertexResponse{Vertex: vertex.AsProto()}, nil
	} else {
		return &pb.GetVertexResponse{}, errors.New("not found")
	}
}

func (l *LanternService) PutVertex(ctx context.Context, request *pb.PutVertexRequest) (*pb.PutVertexResponse, error) {
	for _, vertex := range request.Vertices {
		l.cache.PutVertex(adapter.NewProtoVertex(vertex))
	}
	optProceed.WithLabelValues("put_vertex", "success").Inc()
	return &pb.PutVertexResponse{Status: pb.Status_STATUS_OK}, nil
}

func (l *LanternService) PutEdge(ctx context.Context, request *pb.PutEdgeRequest) (*pb.PutEdgeResponse, error) {
	for _, edge := range request.Edges {
		l.cache.PutEdge(adapter.NewProtoEdge(edge))
	}
	optProceed.WithLabelValues("put_edge", "success").Inc()
	return &pb.PutEdgeResponse{Status: pb.Status_STATUS_OK}, nil
}

func (l *LanternService) Put(ctx context.Context, request *pb.PutRequest) (*pb.PutResponse, error) {
	for _, vertex := range request.Graph.Vertices {
		l.cache.PutVertex(adapter.NewProtoVertex(vertex))
	}
	for _, edge := range request.Graph.Edges {
		l.cache.PutEdge(adapter.NewProtoEdge(edge))
	}

	optProceed.WithLabelValues("put", "success").Inc()
	return &pb.PutResponse{Status: pb.Status_STATUS_OK}, nil
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
	pb.RegisterLanternServiceServer(s.server, s.svc)
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
