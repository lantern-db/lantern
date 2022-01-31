package adapter

import (
	. "github.com/lantern-db/lantern/graph/model"
	pb "github.com/lantern-db/lantern/pb"
)

func LanternQuery(request *pb.IlluminateRequest) LoadQuery {
	return LoadQuery{
		Seed:      Key(request.Seed),
		Step:      request.Step,
		MinWeight: request.MinWeight,
		MaxWeight: request.MaxWeight,
	}
}
