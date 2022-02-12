
./pb/data.pb.go ./pb/data_grpc.pb.go: ./proto/data.proto
	protoc data.proto \
		--proto_path=./proto \
		--go_out=./pb \
		--go-grpc_out=./pb \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative

all: ./pb/data.pb.go ./pb/data_grpc.pb.go ./graph/model/mock/edge.go ./graph/model/mock/vertex.go
.PHONY: all

./graph/model/mock/vertex.go: ./graph/model/vertex.go
	mockgen -source=./graph/model/vertex.go -destination=./graph/model/mock/vertex.go

./graph/model/mock/edge.go: ./graph/model/edge.go
	mockgen -source=./graph/model/edge.go -destination=./graph/model/mock/edge.go

./server/cmd/wire_gen.go: ./server/cmd/wire.go
	wire ./server/cmd/

.PHONY: test
test: ./pb/data.pb.go ./pb/data_grpc.pb.go ./graph/model/mock/vertex.go ./graph/model/mock/edge.go ./server/cmd/wire_gen.go
	go build ./...
	go test ./...

.PHONY: build
build: ./Dockerfile test
	docker build -t lantern .
	docker tag lantern piroyoung/lantern-server:local

.PHONY: run
run: build
	docker run -it -p 6380:6380 piroyoung/lantern-server:local


clean:
	rm ./pb/data.pb.go
	rm ./pb/data_grpc.pb.go
	rm ./graph/model/mock/*
.PHONY: clean