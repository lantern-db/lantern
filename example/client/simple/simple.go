package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/piroyoung/lanterne/client"
	"log"
)

func main() {
	c, err := client.NewLanterneClient("localhost", 6380)
	if err != nil {
		fmt.Printf("hoge %v", err)
		panic(err)
	}
	defer func() {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	_ = c.DumpEdge(ctx, "a", "b", 1.0)
	_ = c.DumpEdge(ctx, "b", "c", 1.0)
	_ = c.DumpEdge(ctx, "c", "d", 1.0)
	_ = c.DumpEdge(ctx, "d", "e", 1.0)

	graph, _ := c.Illuminate(ctx, "a", 2)
	m := jsonpb.Marshaler{}
	jsonString, _ := m.MarshalToString(graph)
	log.Println(jsonString)
}
