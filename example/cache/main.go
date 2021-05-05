package main

import (
	"fmt"
	"github.com/lanternedb/lanterne/graph/cache"
	"github.com/lanternedb/lanterne/graph/model"
	"time"
)

type StringVertex struct {
	value string
}

func (s *StringVertex) Key() string {
	return s.value
}

func (s *StringVertex) Value() interface{} {
	return s.value
}

func main() {
	a := &StringVertex{"a"}
	b := &StringVertex{"b"}
	c := &StringVertex{"c"}
	d := &StringVertex{"d"}
	e := &StringVertex{"e"}

	ab := model.Edge{Tail: a, Head: b, Weight: 1.0}
	bc := model.Edge{Tail: b, Head: c, Weight: 1.0}
	cd := model.Edge{Tail: c, Head: d, Weight: 1.0}
	de := model.Edge{Tail: d, Head: e, Weight: 1.0}
	ce := model.Edge{Tail: c, Head: e, Weight: 1.0}

	repo := cache.NewGraphCache(1 * time.Minute)
	repo.DumpEdge(ab)
	repo.DumpEdge(bc)
	repo.DumpEdge(cd)
	repo.DumpEdge(de)
	repo.DumpEdge(ce)

	q := model.NeighborQuery(a, 5)

	res := repo.Load(q)
	fmt.Println(res.Vertices())
	fmt.Println(res.Edges())

}
