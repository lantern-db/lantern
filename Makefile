.PHONY: init_buf
init_buf:
	cd ./proto && buf mod init

./gen/proto/go/lantern/v1/lantern.pb.go ./gen/proto/go/lantern/v1/lantern_grpc.pb.go: ./proto/lantern/v1/lantern.proto ./buf.gen.yaml
	buf generate proto

all: ./gen/proto/go/lantern/v1/lantern.pb.go ./gen/proto/go/lantern/v1/lantern_grpc.pb.go ./graph/model/mock/edge.go ./graph/model/mock/vertex.go
.PHONY: all

./graph/model/mock/vertex.go: ./graph/model/vertex.go
	mockgen -source=./graph/model/vertex.go -destination=./graph/model/mock/vertex.go

./graph/model/mock/edge.go: ./graph/model/edge.go
	mockgen -source=./graph/model/edge.go -destination=./graph/model/mock/edge.go

./server/cmd/wire_gen.go: ./server/cmd/wire.go
	wire ./server/cmd/

./gateway/cmd/wire_gen.go: ./gateway/cmd/wire.go
	wire ./gateway/cmd/

.PHONY: test
test: ./gen/proto/go/lantern/v1/lantern.pb.go ./gen/proto/go/lantern/v1/lantern_grpc.pb.go ./graph/model/mock/vertex.go ./graph/model/mock/edge.go ./server/cmd/wire_gen.go ./gateway/cmd/wire_gen.go
	go build ./...
	go test ./...

.PHONY: build
build: server.Dockerfile test
	docker build -t lantern-server -f server.Dockerfile .
	docker tag lantern-server piroyoung/lantern-server:local
	docker tag lantern-server piroyoung/lantern-server:latest-alpha
	docker build -t lantern-gateway -f gateway.Dockerfile .
	docker tag lantern-gateway piroyoung/lantern-gateway:local
	docker tag lantern-gateway piroyoung/lantern-gateway:latest-alpha

.PHONY: push
push: build
	docker push piroyoung/lantern-server:latest-alpha
	docker push piroyoung/lantern-gateway:latest-alpha

.PHONY: run
run: build
	docker-compose up


clean:
	rm ./pb/data.pb.go
	rm ./pb/data_grpc.pb.go
	rm ./graph/model/mock/*
.PHONY: clean