
./pb/data.pb.go ./pb/data_grpc.pb.go: ./proto/data.proto
	protoc data.proto \
		--proto_path=./proto \
		--go_out=./pb \
		--go-grpc_out=./pb \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative

all: ./pb/data.pb.go ./pb/data_grpc.pb.go ./graph/model/mock/edge.go ./graph/model/mock/vertex.go
.PHONY: all

./graph/model/mock/vertex.go:
	mockgen -source=./graph/model/vertex.go -destination=./graph/model/mock/vertex.go

./graph/model/mock/edge.go:
	mockgen -source=./graph/model/edge.go -destination=./graph/model/mock/edge.go

clean:
	rm ./pb/data.pb.go
	rm ./pb/data_grpc.pb.go
	rm ./graph/model/mock/*
.PHONY: clean

build: ./Dockerfile ./pb/data.pb.go ./pb/data_grpc.pb.go ./graph/model/mock/vertex.go ./graph/model/mock/edge.go
	docker build -t lantern .
	docker tag lantern piroyoung/lantern-server:latest
	docker tag lantern piroyoung/lantern-server:0.0.4
.PHONY: build
