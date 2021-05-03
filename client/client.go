package client

import (
	"context"
	"errors"
	pb "github.com/piroyoung/lanterne/grpc"
	"google.golang.org/grpc"
	"log"
	"math"
	"strconv"
)

type LanterneClient struct {
	conn   *grpc.ClientConn
	client pb.LanterneClient
}

func New(hostname string, port int) LanterneClient {
	conn, err := grpc.Dial(hostname+":"+strconv.Itoa(port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return LanterneClient{
		conn:   conn,
		client: pb.NewLanterneClient(conn),
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

func (c *LanterneClient) Illuminate(ctx context.Context, seed string, degree uint32) (*pb.Graph, error) {
	request := &pb.IlluminateRequest{
		Seed:      &pb.Vertex{Key: seed},
		Degree:    degree,
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
