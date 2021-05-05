# lantern

[
![DSC00732](https://user-images.githubusercontent.com/6128022/116864177-6824e700-ac42-11eb-8475-c2d06d1761c6.jpg)
](url)

# LanterneDB: 
In recent years, many applications, recommender, fraud detection, are based on a graph structure. And these applications have got more real-time, and dynamic. There are so many graph-based databases, but almost all of graph DB seems too heavy, or too huge.

We've just needed a simple graph structure, but not highly theorized algorithms such as ontologies or optimization techniques.

LanterneDB is In-memory `key-vertex-store` (KVS). 
It behaves like `key-value-store` but can explore neighbor vertices based on graph structure.

LanterneDB is a streaming database.
All vertices or edges will be expired as time passes, just like real relationships.

LanterneDB just illuminates the moment, just focuses on neighbors, not global structures.

# lanterne-server

```
$ docker run -it -p 6380:6380 -e LANTERNE_PORT=6380 -e LANTERNE_TTL=300 piroyoung/lanterne
```

* `LANTERNE_PORT`: Port number for lanterne-server
* `LANTERNE_TTL`: time-to-live for each elements (seconds).

# lanterne-client (Golang)

Example usage of `lanterne-client` for Golang.

`example/client/simple/simple.go`

```golang
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lanternedb/lanterne/client"
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

	_ = c.DumpVertex(ctx, "a", "test")
	_ = c.DumpVertex(ctx, "b", 42)
	_ = c.DumpVertex(ctx, "c", 3.14)

	if resA, err := c.LoadVertex(ctx, "a"); err == nil {
		log.Println(resA.String())
	}
	if resB, err := c.LoadVertex(ctx, "b"); err == nil {
		log.Println(resB.Int())
	}
	if resC, err := c.LoadVertex(ctx, "c"); err == nil {
		log.Println(resC.Float64())
	}
	if resD, err := c.LoadVertex(ctx, "d"); err == nil {
		log.Println(resD.Nil())
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
```

Then we got

```json
{
  "vertexMap": {
    "a": {
      "message": {
        "key": "a",
        "Value": {
          "String_": "test"
        }
      }
    },
    "b": {
      "message": {
        "key": "b",
        "Value": {
          "Int32": 42
        }
      }
    },
    "c": {
      "message": {
        "key": "c",
        "Value": {
          "Float64": 3.14
        }
      }
    }
  },
  "neighborMap": {
    "a": {
      "b": 1
    },
    "b": {
      "c": 1
    }
  }
}
```
