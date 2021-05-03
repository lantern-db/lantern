
lanterne/model/mock.go: graph
	mockgen -source lanterne/model/graph.go -destination lanterne/model/mock.go -package model


./grpc/data.pb.go ./grpc/data_grpc.pb.go: ./proto/data.proto
	protoc data.proto \
		--proto_path=./proto \
		--go_out=./grpc \
		--go-grpc_out=./grpc \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative

all: ./grpc/data.pb.go ./grpc/data_grpc.pb.go
.PHONY: all

clean:
	rm ./grpc/data.pb.go
	rm ./grpc/data_grpc.pb.go
.PHONY: clean