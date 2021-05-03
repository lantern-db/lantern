package main

import (
	"context"
	"fmt"
	"github.com/piroyoung/lanterne/client"
)

func main() {
	c := client.New("localhost", 6380)
	defer c.Close()
	ctx := context.Background()

	_ = c.DumpEdge(ctx, "a", "b", 1.0)
	_ = c.DumpEdge(ctx, "b", "c", 1.0)
	_ = c.DumpEdge(ctx, "c", "d", 1.0)
	_ = c.DumpEdge(ctx, "d", "e", 1.0)

	graph, _ := c.Illuminate(ctx, "a", 2)
	fmt.Println(graph)
	// => Vertices:{key:"a"} Vertices:{key:"b"} Vertices:{key:"c"} Edges:{tail:{key:"a"} head:{key:"b"} weight:1} Edges:{tail:{key:"b"} head:{key:"c"} weight:1}

}
