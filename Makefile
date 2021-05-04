
./graph/model/mock.go: ./graph/model/graph.go
	mockgen -source graph/model/graph.go -destination graph/model/mock.go -package model


./grpc/data.pb.go ./grpc/data_grpc.pb.go: ./proto/data.proto
	protoc data.proto \
		--proto_path=./proto \
		--go_out=./grpc \
		--go-grpc_out=./grpc \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative

all: ./grpc/data.pb.go ./grpc/data_grpc.pb.go ./graph/model/mock.go
.PHONY: all

clean:
	rm ./grpc/data.pb.go
	rm ./grpc/data_grpc.pb.go
.PHONY: clean

build: ./Dockerfile ./graph/model/mock.go ./grpc/data_grpc.pb.go ./grpc/data.pb.go
	docker build -t lanterne .
	docker tag lanterne piroyoung/lanterne:latest
	docker tag lanterne piroyoung/lanterne:0.0.0
	docker push piroyoung/lanterne:latest
	docker push piroyoung/lanterne:0.0.0
.PHONY: build
