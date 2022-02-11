package main

import (
	"fmt"
	"github.com/lantern-db/lantern/graph/adapter"
	"github.com/lantern-db/lantern/graph/cache"
	"github.com/lantern-db/lantern/graph/model"
	"time"
)

func main() {
	ab := adapter.NewProtoEdgeOf("a", "b", 1.0, 60*time.Second)
	bc := adapter.NewProtoEdgeOf("b", "c", 1.0, 60*time.Second)
	cd := adapter.NewProtoEdgeOf("c", "d", 1.0, 60*time.Second)
	de := adapter.NewProtoEdgeOf("d", "e", 1.0, 60*time.Second)
	ce := adapter.NewProtoEdgeOf("c", "e", 1.0, 60*time.Second)

	repo := cache.NewEmptyGraphCache()
	repo.DumpEdge(ab)
	repo.DumpEdge(bc)
	repo.DumpEdge(cd)
	repo.DumpEdge(de)
	repo.DumpEdge(ce)

	q := model.NeighborQuery("a", 3)

	res := repo.Load(q)
	fmt.Println(res.Vertices())
	fmt.Println(res.Edges())

}
