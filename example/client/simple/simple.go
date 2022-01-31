package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lantern-db/lantern/client"
	"log"
)

func main() {
	c, err := client.NewLanternClient("localhost", 6380)
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

	_ = c.DumpVertex(ctx, "a", "test")
	_ = c.DumpVertex(ctx, "b", 42)
	_ = c.DumpVertex(ctx, "c", 3.14)

	if resA, err := c.LoadVertex(ctx, "a"); err == nil {
		log.Println(resA)
		log.Println(resA.StringValue())
	}
	if resB, err := c.LoadVertex(ctx, "b"); err == nil {
		log.Println(resB)
		log.Println(resB.IntValue())
	}
	if resC, err := c.LoadVertex(ctx, "c"); err == nil {
		log.Println(resC)
		log.Println(resC.Float64Value())
	}
	if resD, err := c.LoadVertex(ctx, "d"); err == nil {
		log.Println(resD)
		log.Println(resD.NilValue())
	}

	_ = c.DumpEdge(ctx, "a", "b", 1.0)
	_ = c.DumpEdge(ctx, "b", "c", 1.0)
	_ = c.DumpEdge(ctx, "c", "d", 1.0)
	_ = c.DumpEdge(ctx, "d", "e", 1.0)

	result, err := c.Illuminate(ctx, "a", 2)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	jsonBytes, _ := json.Marshal(result)
	log.Println(string(jsonBytes))
}
