# lanterne

[
![DSC00732](https://user-images.githubusercontent.com/6128022/116864177-6824e700-ac42-11eb-8475-c2d06d1761c6.jpg)
](url)

Lanterns illuminate just your neighbors. Lanterns light only this moment up.

Most of the relations will disappear as time passes. In the case of treating something like social networks, the elapsed
time is an important feature to understand these structures.

`Lanterne` is an in-memory, graph-based, streaming database. Each element like `Edge` or `Vertex` has `time to live`,
and disappears as time passes just like real relationships.

# lanterne-server

```
$ docker run -it -p 6380:6380 -e LANTERNE_PORT=6380 -e LANTERNE_TTL=300 piroyoung/lanterne
```

* `LANTERN_PORT`: Port number for lanterne-server
* `LANTERN_TTL`: time-to-live for each elements (seconds).

# lanterne-client (Golang)

Example usage of `lanterne-client` for Golang.

`example/client/simple/simple.go`

```golang
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

```

Then we got

```json
{
  "vertices": [
    {
      "key": "a"
    },
    {
      "key": "b"
    },
    {
      "key": "c"
    }
  ],
  "edges": [
    {
      "tail": {
        "key": "a"
      },
      "head": {
        "key": "b"
      },
      "weight": 1
    },
    {
      "tail": {
        "key": "b"
      },
      "head": {
        "key": "c"
      },
      "weight": 1
    }
  ]
}
```