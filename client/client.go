package client

import (
	"context"
	"errors"
	pb "github.com/piroyoung/lanterne/grpc"
	"google.golang.org/grpc"
	"math"
	"strconv"
	"time"
)

type LanterneClient struct {
	conn   *grpc.ClientConn
	client pb.LanterneClient
}

func NewLanterneClient(hostname string, port int) (*LanterneClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	chConn := make(chan *grpc.ClientConn)
	chErr := make(chan error)

	go func() {
		conn, err := grpc.DialContext(ctx, hostname+":"+strconv.Itoa(port), grpc.WithInsecure())
		if err != nil {
			chErr <- err
		}
		chConn <- conn
	}()
	select {
	case <-ctx.Done():
		return nil, errors.New("grpc connection timeout")

	case err := <-chErr:
		return nil, err

	case conn := <-chConn:
		return &LanterneClient{
			conn:   conn,
			client: pb.NewLanterneClient(conn),
		}, nil
	}
}

func (c *LanterneClient) Close() error {
	return c.conn.Close()
}

func (c *LanterneClient) DumpEdge(ctx context.Context, tail string, head string, weight float32) error {
	edge := &pb.Edge{
		Tail:   &pb.Vertex{Key: tail},
		Head:   &pb.Vertex{Key: head},
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

func (c *LanterneClient) DumpVertex(ctx context.Context, key string) error {
	response, err := c.client.DumpVertex(ctx, &pb.Vertex{Key: key})
	if err != nil {
		return err
	}
	if response.Status != pb.Status_OK {
		return errors.New("dump vertex error")
	}
	return nil
}

func (c *LanterneClient) Illuminate(ctx context.Context, seed string, step uint32) (*pb.Graph, error) {
	request := &pb.IlluminateRequest{
		Seed:      &pb.Vertex{Key: seed},
		Step:      step,
		MinWeight: -math.MaxFloat32,
		MaxWeight: math.MaxFloat32,
	}
	response, err := c.client.Illuminate(ctx, request)
	if err != nil {
		return nil, err
	}
	if response.Status != pb.Status_OK {
		return nil, errors.New("illuminate error")
	}
	return response.Graph, nil
}
