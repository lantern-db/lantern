# lanterne
[
![DSC00732](https://user-images.githubusercontent.com/6128022/116864177-6824e700-ac42-11eb-8475-c2d06d1761c6.jpg)
](url)

Lanterns illuminate just your neighbors. Lanterns light only this moment up. 

Most of the relations will disappear as time passes. In the case of treating something like social networks, the elapsed time is an important feature to understand these structures.

`Lanterne` is an in-memory, graph-based, streaming database. Each element like `Edge` or `Vertex` has `time to live`, and disappears as time passes just like real relationships.

# lanterne-server

```
$ docker run -it -p 6380:6380 -e LANTERNE_PORT=6380 -e LANTERNE_TTL=300 piroyoung/lanterne
```

`LANTERN_PORT`: Port number for `lanterne-server`
`LANTERN_TTL`: `time-to-live` for each elements (seconds).

# lanterne-client (Golang)
Example usage of `lanterne-client` for Golang.

```golang
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
```
