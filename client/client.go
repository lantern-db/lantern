package client

import (
	"context"
	"errors"
	pb "github.com/lantern-db/lantern/gen/proto/go/lantern/v1"
	"github.com/lantern-db/lantern/graph/adapter"
	"github.com/lantern-db/lantern/graph/model"
	"google.golang.org/grpc"
	"math"
	"strconv"
	"time"
)

type LanternClient struct {
	conn   *grpc.ClientConn
	client pb.LanternServiceClient
}

func NewLanternClient(hostname string, port int) (*LanternClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	chConn := make(chan *grpc.ClientConn)
	chErr := make(chan error)

	go func() {
		conn, err := grpc.DialContext(ctx, hostname+":"+strconv.Itoa(port), grpc.WithInsecure())
		if err != nil {
			chErr <- err
			return
		}
		chConn <- conn
	}()
	select {
	case <-ctx.Done():
		return nil, errors.New("grpc connection timeout")

	case err := <-chErr:
		return nil, err

	case conn := <-chConn:
		return &LanternClient{
			conn:   conn,
			client: pb.NewLanternServiceClient(conn),
		}, nil
	}
}

func (c *LanternClient) Close() error {
	return c.conn.Close()
}

func (c *LanternClient) DumpEdge(ctx context.Context, tail string, head string, weight float32, ttl time.Duration) error {
	edge := &pb.Edge{
		Tail:       tail,
		Head:       head,
		Weight:     weight,
		Expiration: model.NewExpiration(ttl).AsProtoTimestamp(),
	}
	response, err := c.client.PutEdge(ctx, &pb.PutEdgeRequest{Edges: []*pb.Edge{edge}})
	if err != nil {
		return err
	}
	if response.Status != pb.Status_STATUS_OK {
		return errors.New("dump edge error")
	}
	return nil
}

func (c *LanternClient) DumpVertex(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	vertex, err := adapter.NewProtoVertexOf(model.Key(key), model.Value(value), ttl)
	if err != nil {
		return err
	}
	response, err := c.client.PutVertex(ctx, &pb.PutVertexRequest{Vertices: []*pb.Vertex{vertex.AsProto()}})
	if err != nil {
		return err
	}
	if response.Status != pb.Status_STATUS_OK {
		return errors.New("dump vertex error. status: " + pb.Status_STATUS_OK.String())
	}
	return nil
}

func (c *LanternClient) LoadVertex(ctx context.Context, key string) (model.Vertex, error) {
	lanternGraph, err := c.Illuminate(ctx, key, 0)
	if err != nil {
		return nil, err
	}

	r, ok := lanternGraph.VertexMap[model.Key(key)]
	if ok {
		return r, nil
	} else {
		return nil, errors.New("missing values")
	}
}

func (c *LanternClient) Illuminate(ctx context.Context, seed string, step uint32) (*model.Graph, error) {
	request := &pb.IlluminateRequest{
		Seed:      seed,
		Step:      step,
		MinWeight: -math.MaxFloat32,
		MaxWeight: math.MaxFloat32,
	}
	response, err := c.client.Illuminate(ctx, request)
	if err != nil {
		return nil, err
	}
	if response.Status != pb.Status_STATUS_OK {
		return nil, errors.New("illuminate error. status: " + response.Status.String())
	}
	lanternGraph := adapter.LanternGraph(response.Graph)
	return &lanternGraph, nil
}
