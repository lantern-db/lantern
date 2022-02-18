# lantern

[
![DSC00732](https://user-images.githubusercontent.com/6128022/116864177-6824e700-ac42-11eb-8475-c2d06d1761c6.jpg)
](url)

# LanternDB: key-vertex-store

In recent years, many applications, recommender, fraud detection, are based on a graph structure. And these applications
have got more real-time, and dynamic. There are so many graph-based databases, but almost all of graph DB seems too
heavy, or too huge.

We've just needed a simple graph structure, but not highly theorized algorithms such as ontologies or optimization
techniques.

LanternDB is In-memory `key-vertex-store` (KVS). It behaves like `key-value-store` but can explore neighbor vertices
based on graph structure.

LanternDB is a streaming database. All vertices or edges will be expired as time passes, just like real relationships.

LanternDB just illuminates the moment, just focuses on neighbors, not global structures.

# lantern-server

```
$ docker run -it -p 6380:6380 -e LANTERN_PORT=6380 -e piroyoung/lantern-server:alpha
```

* `LANTERN_PORT`: Port number for Lantern-server
* `LANTERN_FLUSH_INTERVAL`: flush interval for expired values

# lantern-client (Golang)

## Setting vertices

We can use `DumpVertex(ctx, key, value, ttl)` method to register vertices like below.

```golang
_ := c.DumpVertex(ctx, "a", "value of a", 60*time.Second)
_ := c.DumpVertex(ctx, "b", "value of b", 60*time.Second)
_ := c.DumpVertex(ctx, "c", "value of c", 60*time.Second)
_ := c.DumpVertex(ctx, "d", "value of d", 60*time.Second)
```

## Loading a single vertex with Key

Loading vertices behaves like standard key-value-store.

```golang
resA, err := c.LoadVertex(ctx, "a")
```

The value of `resA` is an instance of `model.Vertex`, and a field `Vertex.Value` is just `interface{}`. To get the
typed-value, we can use `Vertex.XXXValue()` method such like `IntValue()`. And if a value is not valid as int, this
method will return error.

```golang
str, err := resA.StringValue() // => valid case, returns "value of a"
i, err := resA.IntValue() // => invalid case, returns error
```

## Setting edges

Edges can also be created with `Dump(ctx, keyOfTail, keyOfHead, weight, ttl)`.

```golang
_ = c.DumpEdge(ctx, "a", "b", 1.0, 60*time.Second)
_ = c.DumpEdge(ctx, "b", "c", 1.0, 60*time.Second)
_ = c.DumpEdge(ctx, "c", "d", 1.0, 60*time.Second)
_ = c.DumpEdge(ctx, "b", "e", 1.0, 60*time.Second)
```

If the vertex which has key `a` is missing in a graph, then empty valued vertices will be created with same expirations.

## Incremental weight
Once you set multiple duplicate edges, this weight of the edge will be incremented.

```golang
_ = c.DumpEDge(ctx, "a", "b", 1.0, 3*time.Second) // weight of e(a, b) -> 1.0
_ = c.DumpEDge(ctx, "a", "b", 1.0, 3*time.Second) // weight of e(a, b) -> 2.0

```

But each TTLs of transactions will be expired independently.

### example
```golang
_ = c.DumpEDge(ctx, "a", "b", 1.0, 2*time.Second) // weight of e(a, b) -> 1.0
time.Sleep(1*time.Second)                         // weight of e(a, b) -> 1.0
_ = c.DumpEDge(ctx, "a", "b", 1.0, 2*time.Second) // weight of e(a, b) -> 2.0
time.Sleep(1*time.Second)                         // weight of e(a, b) -> 1.0, first transaction is expired
time.Sleep(1*time.Second)                         // weight of e(a, b) -> 0.0, second transaction is expired
```


## Loading vertices and its neighbors with key and n_step

All vertices can be loaded with a graph structure linking with edges, And we call this transaction `Illuminate`. An
client has method `Illuminate(ctx, key, step)`.

For example, `client.Illuminate(ctx, "a", 2)` returns all vertices within 2-steps from a vertex "a". type of returning
value is an instance of `model.Graph` and it can be rendered to json-parsable struct with a method `Render()`.

```json
{
  "vertices": {
    "a": "test",
    "b": 42,
    "c": 3.14,
    "d": null,
    "e": null
  },
  "edges": {
    "a": {
      "b": 1
    },
    "b": {
      "c": 1,
      "e": 1
    },
    "c": {
      "d": 1
    }
  }
}
```

## Brief Example

Example usage of `Lantern-client` for Golang.

`example/client/simple/simple.go`

```golang

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lantern-db/lantern/client"
	"log"
	"os"
	"time"
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

	_ = c.DumpVertex(ctx, "a", "test", 60*time.Second)
	_ = c.DumpVertex(ctx, "b", 42, 60*time.Second)
	_ = c.DumpVertex(ctx, "c", 3.14, 60*time.Second)

	if resA, err := c.LoadVertex(ctx, "a"); err == nil {
		log.Println(resA.StringValue())
	}
	if resB, err := c.LoadVertex(ctx, "b"); err == nil {
		log.Println(resB.IntValue())
	}
	if resC, err := c.LoadVertex(ctx, "c"); err == nil {
		log.Println(resC.Float64Value())
	}
	if _, err := c.LoadVertex(ctx, "d"); err == nil {
		log.Println(err)
	}

	_ = c.DumpEdge(ctx, "a", "b", 1.0, 60*time.Second)
	_ = c.DumpEdge(ctx, "b", "c", 1.0, 60*time.Second)
	_ = c.DumpEdge(ctx, "c", "d", 1.0, 60*time.Second)
	_ = c.DumpEdge(ctx, "b", "e", 1.0, 60*time.Second)

	result, err := c.Illuminate(ctx, "a", 3)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	jsonBytes, _ := json.Marshal(result.Render())
	log.Println(string(jsonBytes))
}
```

Then we will get

```json
{
  "vertices": {
    "a": "test",
    "b": 42,
    "c": 3.14,
    "d": null,
    "e": null
  },
  "edges": {
    "a": {
      "b": 1
    },
    "b": {
      "c": 1,
      "e": 1
    },
    "c": {
      "d": 1
    }
  }
}
```

And it shows simple graph structure.


