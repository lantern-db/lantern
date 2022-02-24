package adapter

import (
	pb "github.com/lantern-db/lantern/gen/proto/go/lantern/v1"
	. "github.com/lantern-db/lantern/graph/model"
)

func LanternQuery(request *pb.IlluminateRequest) LoadQuery {
	return LoadQuery{
		Seed: Key(request.Seed),
		Step: request.Step,
		TopK: request.TopK,
	}
}
