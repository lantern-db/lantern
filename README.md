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
$ docker run -it -p 6380:6380 -e LANTERN_PORT=6380 -e LANTERN_TTL=300 piroyoung/lantern-server
```

* `LANTERN_PORT`: Port number for Lantern-server
* `LANTERN_TTL`: time-to-live for each elements (seconds).

# lantern-client (Golang)

## Setting vertices

We can use `DumpVertex(ctx, key, value)` method to register vertices like below.

```golang
_ := c.DumpVertex(ctx, "a", "value of a")
_ := c.DumpVertex(ctx, "b", "value of b")
_ := c.DumpVertex(ctx, "c", "value of c")
_ := c.DumpVertex(ctx, "d", "value of d")
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

Edges can also be created with `Dump(ctx, keyOfTail, keyOfHead, weight)`.

```golang
_ = c.DumpEdge(ctx, "a", "b", 1.0)
_ = c.DumpEdge(ctx, "b", "c", 1.0)
_ = c.DumpEdge(ctx, "c", "d", 1.0)
_ = c.DumpEdge(ctx, "d", "e", 1.0)
```

If the vertex which has key `a` is missing in a graph, then empty valued vertices will be created with same expirations.

## Loading vertices and its neighbors with key and n_step

All vertices can be loaded with a graph structure linking with edges, And we call this transaction `Illuminate`. An
client has method `Illuminate(ctx, key, step)`.

For Example, `client.Illuminate(ctx, "a", 2)` returns all vertices within 2-steps from a vertex "a". type of returning
value is an instance of `model.Graph` and it's json-parsable struct.

```json
{
  "vertexMap": {
    "a": {
      "key": "a",
      "value": "value of a",
      "expiration": 1643642189
    },
    "b": {
      "key": "b",
      "value": "value of b",
      "expiration": 1643642189
    },
    "c": {
      "key": "c",
      "value": "value of c",
      "expiration": 1643642189
    }
  },
  "edgeMap": {
    "a": {
      "b": {
        "tail": "a",
        "head": "b",
        "weight": 1,
        "expiration": 1643642189
      }
    },
    "b": {
      "c": {
        "tail": "b",
        "head": "c",
        "weight": 1,
        "expiration": 1643642189
      }
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
		log.Println(resA.StringValue())
	}
	if resB, err := c.LoadVertex(ctx, "b"); err == nil {
		log.Println(resB.IntValue())
	}
	if resC, err := c.LoadVertex(ctx, "c"); err == nil {
		log.Println(resC.Float64Value())
	}
	if resD, err := c.LoadVertex(ctx, "d"); err == nil {
		log.Println(resD.NilValue())
	}

	_ = c.DumpEdge(ctx, "a", "b", 1.0)
	_ = c.DumpEdge(ctx, "b", "c", 1.0)
	_ = c.DumpEdge(ctx, "c", "d", 1.0)
	_ = c.DumpEdge(ctx, "d", "e", 1.0)

	result, err := c.Illuminate(ctx, "a", 3)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	jsonBytes, _ := json.Marshal(result)
	log.Println(string(jsonBytes))
}
```

Then we got

```json
{
  "vertexMap": {
    "a": {
      "key": "a",
      "value": "test",
      "expiration": 1643642189
    },
    "b": {
      "key": "b",
      "value": 42,
      "expiration": 1643642189
    },
    "c": {
      "key": "c",
      "value": 3.14,
      "expiration": 1643642189
    },
    "d": {
      "key": "d",
      "expiration": 1643642040
    }
  },
  "edgeMap": {
    "a": {
      "b": {
        "tail": "a",
        "head": "b",
        "weight": 1,
        "expiration": 1643642189
      }
    },
    "b": {
      "c": {
        "tail": "b",
        "head": "c",
        "weight": 1,
        "expiration": 1643642189
      }
    },
    "c": {
      "d": {
        "tail": "c",
        "head": "d",
        "weight": 1,
        "expiration": 1643642189
      }
    }
  }
}
```

And it shows simple graph structure.


