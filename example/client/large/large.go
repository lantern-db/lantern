package main

import (
	"context"
	"encoding/json"
	"github.com/lantern-db/lantern/client"
	"log"
	"math/rand"
	"strconv"
)

func main() {
	c, err := client.NewLanternClient("localhost", 6380)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	var i, j int
	i = 0
	for range make([]int, 1000) {
		j = rand.Intn(100)
		log.Println(i, j)
		err := c.DumpEdge(ctx, strconv.Itoa(i), strconv.Itoa(j), 1.0)
		if err != nil {
			panic(err)
		}
		i = j
	}

	result, err := c.Illuminate(ctx, "0", 5)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	jsonBytes, _ := json.Marshal(result)
	log.Println(string(jsonBytes))
}
