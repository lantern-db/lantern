package main

import (
	"context"
	"encoding/json"
	"github.com/lantern-db/lantern/client"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
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
		err = c.DumpVertex(ctx, strconv.Itoa(i), strconv.Itoa(i), 60*time.Minute)
		err = c.DumpVertex(ctx, strconv.Itoa(j), strconv.Itoa(j), 60*time.Minute)
		err = c.DumpEdge(ctx, strconv.Itoa(i), strconv.Itoa(j), 1.0, 60*time.Minute)
		if err != nil {
			panic(err)
		}
		i = j
	}

	result, err := c.Illuminate(ctx, "0", 5)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	jsonBytes, _ := json.Marshal(result.Render())
	log.Println(string(jsonBytes))
	if err := writeBytes("./large.json", string(jsonBytes)); err != nil {
		log.Fatal(err)
	}
}

func writeBytes(filename string, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	b := []byte(content)
	_, err = file.Write(b)
	if err != nil {
		return err
	}
	return nil
}
