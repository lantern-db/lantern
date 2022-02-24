package main

import (
	"encoding/json"
	"github.com/lantern-db/lantern/graph/adapter"
	"github.com/lantern-db/lantern/graph/cache"
	"github.com/lantern-db/lantern/graph/model"
	"log"
	"time"
)

func main() {
	ab := adapter.NewProtoEdgeOf("a", "b", 1.0, 60*time.Second)
	bc := adapter.NewProtoEdgeOf("b", "c", 1.0, 60*time.Second)
	cd := adapter.NewProtoEdgeOf("c", "d", 1.0, 60*time.Second)
	de := adapter.NewProtoEdgeOf("d", "e", 1.0, 60*time.Second)
	ce := adapter.NewProtoEdgeOf("c", "e", 1.0, 60*time.Second)

	repo := cache.NewEmptyGraphCache()
	repo.PutEdge(ab)
	repo.PutEdge(bc)
	repo.PutEdge(cd)
	repo.PutEdge(de)
	repo.PutEdge(ce)
	v, _ := adapter.NewProtoVertexOf("a", 0, 60*time.Second)
	repo.PutVertex(v)

	q := model.NeighborQuery("a", 4, 10)

	res := repo.Load(q)
	jsonBytes, _ := json.Marshal(res.Render())
	log.Println(string(jsonBytes))

}
