package client

import (
	"context"
	"errors"
	"github.com/lantern-db/lantern/graph/adapter"
	"github.com/lantern-db/lantern/graph/model"
	pb "github.com/lantern-db/lantern/pb"
	"google.golang.org/grpc"
	"math"
	"strconv"
	"time"
)

type LanternClient struct {
	conn   *grpc.ClientConn
	client pb.LanternClient
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
			client: pb.NewLanternClient(conn),
		}, nil
	}
}

func (c *LanternClient) Close() error {
	return c.conn.Close()
}

func (c *LanternClient) DumpEdge(ctx context.Context, tail string, head string, weight float32) error {
	edge := &pb.Edge{
		Tail:   tail,
		Head:   head,
		Weight: weight,
	}
	response, err := c.client.DumpEdge(ctx, edge)
	if err != nil {
		return err
	}
	if response.Status != pb.Status_OK {
		return errors.New("dump edge error")
	}
	return nil
}

func (c *LanternClient) DumpVertex(ctx context.Context, key string, value interface{}) error {
	vertex := adapter.ProtoVertex(model.Vertex{Key: model.Key(key), Value: value})
	response, err := c.client.DumpVertex(ctx, vertex)
	if err != nil {
		return err
	}
	if response.Status != pb.Status_OK {
		return errors.New("dump vertex error. status: " + pb.Status_OK.String())
	}
	return nil
}

func (c *LanternClient) LoadVertex(ctx context.Context, key string) (*model.Vertex, error) {
	lanternGraph, err := c.Illuminate(ctx, key, 0)
	if err != nil {
		return nil, err
	}
	r := lanternGraph.VertexMap[model.Key(key)]
	return &r, nil
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
	if response.Status != pb.Status_OK {
		return nil, errors.New("illuminate error. status: " + response.Status.String())
	}
	lanternGraph := adapter.LanternGraph(response.Graph)
	return &lanternGraph, nil
}
