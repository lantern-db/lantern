
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

build: ./Dockerfile ./pb/data.pb.go ./pb/data_grpc.pb.go
	docker build -t lantern .
	docker tag lantern piroyoung/lantern-server:latest
	docker tag lantern piroyoung/lantern-server:0.0.4
	docker push piroyoung/lantern-server:latest
	docker push piroyoung/lantern-server:0.0.4
.PHONY: build
