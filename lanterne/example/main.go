package main

import (
	"context"
	"fmt"
	"github.com/piroyoung/lanterne/lanterne/model"
	"github.com/piroyoung/lanterne/lanterne/repository"
	"time"
)

type StringVertex struct {
	value string
}

func (s StringVertex) Digest() string {
	return s.value
}

func main() {
	a := &StringVertex{"a"}
	b := &StringVertex{"b"}
	c := &StringVertex{"c"}
	d := &StringVertex{"d"}
	e := &StringVertex{"e"}

	ab := model.Edge{a, b, 1.0}
	bc := model.Edge{b, c, 1.0}
	cd := model.Edge{c, d, 1.0}
	de := model.Edge{d, e, 1.0}
	ce := model.Edge{c, e, 1.0}

	ctx := context.Background()
	repo := repository.NewCacheGraphRepository(1 * time.Minute)
	_ = repo.DumpEdge(ctx, ab)
	_ = repo.DumpEdge(ctx, bc)
	_ = repo.DumpEdge(ctx, cd)
	_ = repo.DumpEdge(ctx, de)
	_ = repo.DumpEdge(ctx, ce)

	q := model.NeighborQuery{
		Seed:      a,
		Degree:    3,
		MinWeight: -1,
		MaxWeight: 100,
	}

	res, err := repo.LoadNeighbor(ctx, q)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

}
