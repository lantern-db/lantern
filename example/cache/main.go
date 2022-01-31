package main

import (
	"fmt"
	"github.com/lantern-db/lantern/graph/cache"
	"github.com/lantern-db/lantern/graph/model"
	"time"
)

func main() {
	ab := model.Edge{Tail: "a", Head: "b", Weight: 1.0}
	bc := model.Edge{Tail: "b", Head: "c", Weight: 1.0}
	cd := model.Edge{Tail: "c", Head: "d", Weight: 1.0}
	de := model.Edge{Tail: "d", Head: "e", Weight: 1.0}
	ce := model.Edge{Tail: "c", Head: "e", Weight: 1.0}

	repo := cache.NewEmptyGraphCache(1 * time.Minute)
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
