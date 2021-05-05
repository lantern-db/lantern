
./graph/model/mock.go: ./graph/model/graph.go
	mockgen -source graph/model/graph.go -destination graph/model/mock.go -package model


./pb/data.pb.go ./pb/data_grpc.pb.go: ./proto/data.proto
	protoc data.proto \
		--proto_path=./proto \
		--go_out=./pb \
		--go-grpc_out=./pb \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative

all: ./pb/data.pb.go ./pb/data_grpc.pb.go ./graph/model/mock.go
.PHONY: all

clean:
	rm ./pb/data.pb.go
	rm ./pb/data_grpc.pb.go
	rm ./graph/model/mock.go
.PHONY: clean

build: ./Dockerfile ./graph/model/mock.go ./pb/data.pb.go ./pb/data_grpc.pb.go
	docker build -t lanterne .
	docker tag lanterne piroyoung/lanterne:latest
	docker tag lanterne piroyoung/lanterne:0.0.0
	docker push piroyoung/lanterne:latest
	docker push piroyoung/lanterne:0.0.0
.PHONY: build
